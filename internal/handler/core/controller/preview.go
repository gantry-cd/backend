package controller

import (
	"context"
	"errors"
	"fmt"
	"github.com/gantrycd/backend/cmd/config"
	coreErr "github.com/gantrycd/backend/internal/error"
	"github.com/gantrycd/backend/internal/usecases/core/k8sclient"
	"github.com/gantrycd/backend/internal/utils"
	"github.com/gantrycd/backend/internal/utils/branch"
	v1 "github.com/gantrycd/backend/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *controller) CreatePreview(ctx context.Context, in *v1.CreatePreviewRequest) (*v1.CreatePreviewReply, error) {
	branchName := branch.Transpile1123(in.Branch)
	dep, err := c.control.GetDeployment(ctx,
		k8sclient.GetDeploymentParams{
			Namespace:     in.Organization,
			Repository:    in.Repository,
			PullRequestID: in.PullRequestId,
			Branch:        branchName,
		})
	if err != nil && !errors.Is(err, coreErr.ErrDeploymentsNotFound) {
		return nil, err
	}

	if dep != nil {
		return &v1.CreatePreviewReply{
			Name:      dep.Name,
			Namespace: dep.Namespace,
			Version:   dep.ResourceVersion,
		}, nil
	}

	return c.createDeployment(ctx, in)
}

func (c *controller) createDeployment(ctx context.Context, in *v1.CreatePreviewRequest) (*v1.CreatePreviewReply, error) {
	branchName := branch.Transpile1123(in.Branch)

	if in.Replicas == 0 {
		in.Replicas = 1
	}

	deps, err := c.control.CreateDeployment(ctx,
		k8sclient.CreateDeploymentParams{
			Namespace: in.Organization,
			AppName:   in.Repository,
			Image:     in.Image,
			Config:    in.Configs,
			Replicas:  in.Replicas,
		},
		k8sclient.WithRepositoryLabel(in.Repository),
		k8sclient.WithPrIDLabel(in.PullRequestId),
		k8sclient.WithEnvirionmentLabel(k8sclient.EnvPreview),
		k8sclient.WithBaseBranchLabel(branchName),
		k8sclient.WithCreatedByLabel(k8sclient.AppIdentifier),
	)
	if err != nil {
		return nil, err
	}

	// NodePortの指定がなかったら終了
	if in.ExposePorts == nil {
		return &v1.CreatePreviewReply{
			Name:      deps.Name,
			Namespace: deps.Namespace,
			Version:   deps.ResourceVersion,
		}, nil
	}

	baseDomain := fmt.Sprintf("%s-%s-%s", in.Organization, in.Repository, in.PullRequestId)

	service, err := c.control.CreateLoadBalancerService(ctx,
		k8sclient.CreateServiceNodePortParams{
			Namespace:   in.Organization,
			ServiceName: deps.Name,
			TargetPort:  in.ExposePorts,
		},
		k8sclient.WithAppLabel(in.Repository),
		k8sclient.WithRepositoryLabel(in.Repository),
		k8sclient.WithPrIDLabel(in.PullRequestId),
		k8sclient.WithEnvirionmentLabel(k8sclient.EnvPreview),
		k8sclient.WithBaseBranchLabel(branchName),
		k8sclient.WithCreatedByLabel(k8sclient.AppIdentifier),
	)
	if err != nil {
		return nil, err
	}

	cloudflaredConfigYaml, domains := buildCloudflaredConfig(in.Organization, service.Name, baseDomain, in.ExposePorts)

	configMapName := fmt.Sprintf("%s-configMap", baseDomain)
	cloudflaredPodName := fmt.Sprintf("%s-cloudflared", baseDomain)

	if err := c.control.CreateConfigMap(ctx, corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: in.Organization,
		},
		Data: map[string]string{
			"config.yaml": cloudflaredConfigYaml,
		},
	}, metav1.CreateOptions{}); err != nil {
		return nil, err
	}

	if _, err := c.control.CreatePod(ctx, in.Organization, &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cloudflaredPodName,
			Namespace: in.Organization,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  cloudflaredPodName,
					Image: "cloudflare/cloudflared:latest",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "config",
							MountPath: "/etc/cloudflared/config/",
						},
					},
					Args: []string{"tunnel", "--config", "/etc/cloudflared/config/config.yaml", "run"},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "config",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: configMapName,
							},
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{}); err != nil {
		return nil, err
	}

	return &v1.CreatePreviewReply{
		Name:      deps.Name,
		Namespace: deps.Namespace,
		Version:   deps.ResourceVersion,
		External:  domains,
	}, nil
}

func buildCloudflaredConfig(namespace string, serviceName string, baseDomain string, ports []int32) (string, []string) {
	var ingress = ""
	var domains []string
	for _, port := range ports {
		domains = append(domains, fmt.Sprintf("%s-%d.%s", baseDomain, port, config.Config.Application.ExternalDomain))
		ingress += fmt.Sprintf(`  - hostname: %s-%d.%s
    service: http://%s.%s.svc.cluster.local:%d
`, baseDomain, port, config.Config.Application.ExternalDomain, serviceName, namespace, port)
	}
	return fmt.Sprintf(`tunnel: %s
credentials-file: /etc/cloudflared/credentials.json
no-autoupdate: true

ingress:
%s`, config.Config.Application.CloudflaredTunnelId, ingress), domains
}

func (c *controller) UpdatePreview(ctx context.Context, in *v1.CreatePreviewRequest) (*v1.CreatePreviewReply, error) {
	branchName := branch.Transpile1123(in.Branch)

	dep, err := c.control.GetDeployment(ctx,
		k8sclient.GetDeploymentParams{
			Namespace:     in.Organization,
			Repository:    in.Repository,
			PullRequestID: in.PullRequestId,
			Branch:        branchName,
		})
	if err != nil && !errors.Is(err, coreErr.ErrDeploymentsNotFound) {
		return nil, err
	}

	if dep == nil {
		return c.createDeployment(ctx, in)
	}

	dep, err = c.control.UpdateDeployment(ctx, dep, k8sclient.UpdateDeploymentParams{
		Namespace:     in.Organization,
		Repository:    in.Repository,
		PullRequestID: in.PullRequestId,
		Branch:        branch.Transpile1123(in.Branch),
		AppName:       in.Repository,
		Image:         utils.ToPtr(in.Image),
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreatePreviewReply{
		Name:      dep.Name,
		Namespace: dep.Namespace,
		Version:   dep.ResourceVersion,
	}, nil
}

func (c *controller) DeletePreview(ctx context.Context, in *v1.DeletePreviewRequest) (*emptypb.Empty, error) {
	branchName := branch.Transpile1123(in.Branch)
	if err := c.control.DeleteDeployment(ctx,
		in.Organization,
		k8sclient.WithAppLabel(in.Repository),
		k8sclient.WithRepositoryLabel(in.Repository),
		k8sclient.WithBaseBranchLabel(branchName),
		k8sclient.WithPrIDLabel(in.PullRequestId),
		k8sclient.WithCreatedByLabel(k8sclient.AppIdentifier),
		k8sclient.WithEnvirionmentLabel(k8sclient.EnvPreview),
	); err != nil {
		return nil, err
	}

	if err := c.control.DeleteService(ctx,
		in.Organization,
		k8sclient.WithAppLabel(in.Repository),
		k8sclient.WithRepositoryLabel(in.Repository),
		k8sclient.WithBaseBranchLabel(branchName),
		k8sclient.WithPrIDLabel(in.PullRequestId),
		k8sclient.WithCreatedByLabel(k8sclient.AppIdentifier),
		k8sclient.WithEnvirionmentLabel(k8sclient.EnvPreview),
	); err != nil {
		return nil, err
	}

	baseDomain := fmt.Sprintf("%s-%s-%s", in.Organization, in.Repository, in.PullRequestId)
	configMapName := fmt.Sprintf("%s-configMap", baseDomain)
	cloudflaredPodName := fmt.Sprintf("%s-cloudflared", baseDomain)

	if err := c.control.DeleteConfigMap(ctx, in.Organization, configMapName, metav1.DeleteOptions{}); err != nil {
		return nil, err
	}
	if err := c.control.DeletePod(ctx, in.Organization, cloudflaredPodName, metav1.DeleteOptions{}); err != nil {
		return nil, err
	}

	return new(emptypb.Empty), nil
}
