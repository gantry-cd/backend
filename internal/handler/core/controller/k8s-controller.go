package controller

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	coreErr "github.com/gantrycd/backend/internal/error"
	"github.com/gantrycd/backend/internal/usecases/core/k8sclient"
	"github.com/gantrycd/backend/internal/usecases/core/resource"
	"github.com/gantrycd/backend/internal/utils/branch"
	v1 "github.com/gantrycd/backend/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/kubernetes"
)

const (
	// NamespaceLabel is the label used to identify the namespace
	CreatedLabel = "created"
	Identity     = "gantrycd"
)

type controller struct {
	v1.UnimplementedK8SCustomControllerServer
	client  *kubernetes.Clientset
	control k8sclient.K8SClient
	metric  resource.Resource
}

func NewController(
	client *kubernetes.Clientset,
	metric resource.Resource,
) v1.K8SCustomControllerServer {
	return &controller{
		client:  client,
		control: k8sclient.New(client),
		metric:  metric,
	}
}

func (c *controller) CreateNamespace(ctx context.Context, in *v1.CreateNamespaceRequest) (*v1.CreateNamespaceReply, error) {
	ns, err := c.control.CreateNamespace(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &v1.CreateNamespaceReply{
		Name: ns.Name,
	}, nil
}

func (c *controller) ListNamespaces(ctx context.Context, in *emptypb.Empty) (*v1.ListNamespacesReply, error) {
	ns, err := c.control.ListNamespaces(ctx)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, n := range ns.Items {
		names = append(names, n.Name)
	}

	return &v1.ListNamespacesReply{
		Names: names,
	}, nil
}

func (c *controller) DeleteNamespace(ctx context.Context, in *v1.DeleteNamespaceRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.control.DeleteNamespace(ctx, in.Name)
}

func (c *controller) ApplyDeployment(ctx context.Context, in *v1.CreateDeploymentRequest) (*v1.CreateDeploymentReply, error) {
	dep, err := c.control.GetDeployment(ctx, k8sclient.GetDeploymentParams{
		Namespace:     in.Namespace,
		Repository:    in.Repository,
		PullRequestID: in.PrNumber,
		Branch:        branch.Transpile1123(in.Branch),
	})
	if err != nil && !errors.Is(err, coreErr.ErrDeploymentsNotFound) {
		return nil, err
	}

	// リソースが既に存在している場合は、更新する
	if dep != nil {
		// TODO: Apply Deployment
		return &v1.CreateDeploymentReply{
			Name:      dep.Name,
			Namespace: dep.Namespace,
			Version:   dep.ResourceVersion,
		}, nil
	}

	// リソースが存在しない場合は、新規作成する
	dep, err = c.control.CreateDeployment(ctx, k8sclient.CreateDeploymentParams{
		Namespace: in.Namespace,
		AppName:   in.AppName,
		Image:     in.Image,
	},
		k8sclient.WithRepositoryLabel(in.Repository),
		k8sclient.WithPrIDLabel(in.PrNumber),
		k8sclient.WithEnvirionmentLabel(k8sclient.EnvPreview),
		k8sclient.WithBaseBranchLabel(branch.Transpile1123(in.Branch)),
	)
	if err != nil {
		return nil, err
	}

	return &v1.CreateDeploymentReply{
		Name:      dep.Name,
		Namespace: dep.Namespace,
		Version:   dep.ResourceVersion,
	}, nil
}

func (c *controller) DeleteDeployment(ctx context.Context, in *v1.DeleteDeploymentRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.control.DeleteDeployment(ctx, in.Namespace,
		k8sclient.WithRepositoryLabel(in.Repository),
		k8sclient.WithPrIDLabel(in.PrNumber),
		// k8sclient.WithBaseBranchLabel(branch.Transpile1123(in.Branch)),
	)
}

func (c *controller) GetOrgRepos(ctx context.Context, in *v1.GetOrgRepoRequest) (*v1.GetOrgReposReply, error) {
	return c.getOrganization(ctx, in.Organization)
}

func (c *controller) GetResource(ctx context.Context, in *v1.GetResourceRequest) (*v1.GetResourceReply, error) {
	resource, err := c.metric.GetLoads(ctx, in.GetOrganization(), in.GetRepository())
	if err != nil {
		return &v1.GetResourceReply{
			Resources: resource,
			IsDisable: false,
		}, nil
	}

	return &v1.GetResourceReply{
		Resources: resource,
		IsDisable: true,
	}, nil
}

func (c *controller) GetAlls(ctx context.Context, in *emptypb.Empty) (*v1.GetAllsReply, error) {
	namespaces, err := c.control.ListNamespaces(ctx, k8sclient.WithCreatedByLabel(k8sclient.AppIdentifier))
	if err != nil {
		return nil, err
	}

	orgs := []*v1.GetOrgReposReply{}

	for _, ns := range namespaces.Items {
		org, err := c.getOrganization(ctx, ns.Name)
		if err != nil {
			return nil, err
		}

		orgs = append(orgs, org)
	}

	return &v1.GetAllsReply{
		OrganizationInfos: orgs,
	}, nil
}

func (c *controller) getOrganization(ctx context.Context, organization string) (*v1.GetOrgReposReply, error) {
	deployments, err := c.control.ListDeployments(ctx, organization)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list deployments: %v", err)
	}

	var (
		repos []*v1.Repository
		apps  []*v1.Application
	)

	for _, d := range deployments.Items {
		pullRequestID, prOk := d.Labels[k8sclient.PullRequestID]
		branchName, brOk := d.Labels[k8sclient.BaseBranchLabel]
		// PR番号とブランチ名が両方ともない場合はAppとして扱う
		sort.Slice(d.Status.Conditions, func(i, j int) bool {
			return d.Status.Conditions[j].LastTransitionTime.Before(&d.Status.Conditions[i].LastTransitionTime)
		})

		if !prOk && !brOk {
			apps = append(apps, &v1.Application{
				Name:    d.Name,
				Status:  string(d.Status.Conditions[0].Type),
				Version: d.Spec.Template.GetResourceVersion(),
				Image:   d.Spec.Template.Spec.Containers[0].Image,
				Age:     d.CreationTimestamp.Format(time.DateTime),
			})
			continue
		}

		branchName, _ = branch.TranspileBranchName(branchName)

		// PR番号かブランチ名のどちらかが場合はRepoとして扱う
		repos = append(repos, &v1.Repository{
			Name:          d.Labels[k8sclient.RepositoryLabel],
			PullRequestId: pullRequestID,
			Branch:        branchName,
		})
	}

	return &v1.GetOrgReposReply{
		Organization: organization,
		Applications: apps,
		Repositories: repos,
	}, nil
}

func (c *controller) GetBranchInfo(ctx context.Context, in *v1.GetBranchInfoRequest) (*v1.GetBranchInfoReply, error) {
	dep, err := c.control.GetDeployment(ctx, k8sclient.GetDeploymentParams{
		Namespace:     in.Organization,
		Repository:    in.Repository,
		PullRequestID: in.PullreqId,
		Branch:        branch.Transpile1123(in.Branch),
	})
	if err != nil && !errors.Is(err, coreErr.ErrDeploymentsNotFound) {
		return nil, err
	}

	rowYaml, err := dep.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal deployment: %w", err)
	}
	yaml := string(rowYaml)
	return &v1.GetBranchInfoReply{
		Yaml: yaml,
	}, nil
}

func (c *controller) GetDeployInfomation(ctx context.Context, in *v1.GetDeployInfomationRequest) (*v1.GetDeployInfomationReply, error) {
	dep, err := c.control.GetDeployment(ctx, k8sclient.GetDeploymentParams{
		Namespace:     in.Organization,
		Repository:    in.Repository,
		PullRequestID: in.PullRequestId,
		Branch:        branch.Transpile1123(in.Branch),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// objectmeta.managefields.fieldv1のフィールドを削除
	for i := range dep.GetObjectMeta().GetManagedFields() {
		dep.ObjectMeta.ManagedFields[i].FieldsV1 = nil
	}

	depYaml, err := yaml.Marshal(dep)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	services, err := c.control.GetServices(ctx, k8sclient.GetServicesParams{
		Namespace:     in.Organization,
		Repository:    in.Repository,
		PullRequestID: in.PullRequestId,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	for _, svc := range services {
		// objectmeta.managefields.fieldv1のフィールドを削除
		for i := range svc.GetObjectMeta().GetManagedFields() {
			svc.ObjectMeta.ManagedFields[i].FieldsV1 = nil
		}

		svcYaml, err := yaml.Marshal(svc)

		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		depYaml = append(depYaml, []byte(fmt.Sprintf("\n---\n%s", svcYaml))...)
	}

	pods, err := c.control.GetPods(ctx, in.Organization, in.Repository)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var respPods []*v1.Pod
	for _, pod := range pods {
		respPods = append(respPods, &v1.Pod{
			Name:   pod.Name,
			Status: string(pod.Status.Phase),
			Age:    pod.CreationTimestamp.Format(time.DateTime),
			Image:  pod.Spec.Containers[0].Image,
		})
	}

	return &v1.GetDeployInfomationReply{
		Namespace: in.GetOrganization(),
		Branch:    dep.Labels[k8sclient.BaseBranchLabel],
		Pods:      respPods,
		Yaml:      string(depYaml),
	}, nil
}
