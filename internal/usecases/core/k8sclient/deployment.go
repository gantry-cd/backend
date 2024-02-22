package k8sclient

import (
	"context"
	"fmt"

	coreErr "github.com/gantrycd/backend/internal/error"
	pbv1 "github.com/gantrycd/backend/proto"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type CreateDeploymentParams struct {
	Namespace string
	AppName   string
	Image     string
	Replicas  int32
	Config    []*pbv1.Config
}

func (k *k8sClient) CreateDeployment(ctx context.Context, in CreateDeploymentParams, opts ...Option) (*appsv1.Deployment, error) {
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	k.l.Info("creating deployment", "namespace", in.Namespace, "appName", in.AppName, "image", in.Image, CreatedByLabel, o.labelSelector[CreatedByLabel])

	o.labelSelector[AppLabel] = in.AppName

	var config []corev1.EnvVar
	for _, c := range in.Config {
		config = append(config, corev1.EnvVar{
			Name:  c.GetName(),
			Value: c.GetValue(),
		})

	}

	return k.client.AppsV1().Deployments(in.Namespace).Create(ctx, &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", in.AppName),
			Labels:       o.labelSelector,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &in.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: o.labelSelector,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: o.labelSelector,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            in.AppName,
							Image:           in.Image,
							ImagePullPolicy: o.containerOption[in.Image].imagePullPolicy,
							Env:             config,
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{})
}

type GetDeploymentParams struct {
	Namespace     string
	Repository    string
	PullRequestID string
	Branch        string
}

func (k *k8sClient) GetDeployment(ctx context.Context, param GetDeploymentParams) (*appsv1.Deployment, error) {
	deps, err := k.client.AppsV1().Deployments(param.Namespace).List(ctx, metav1.ListOptions{})

	for _, dep := range deps.Items {
		if dep.Labels[RepositoryLabel] == param.Repository && dep.Labels[PullRequestID] == param.PullRequestID && dep.Labels[BaseBranchLabel] == param.Branch {
			return &dep, nil
		}
	}

	if err != nil {
		return nil, err
	}

	return nil, coreErr.ErrDeploymentsNotFound
}

func (k *k8sClient) ListDeployments(ctx context.Context, namespace string, opts ...Option) (*appsv1.DeploymentList, error) {
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	return k.client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{LabelSelector: labels.Set(o.labelSelector).String()})
}

type UpdateDeploymentParams struct {
	Namespace     string
	Repository    string
	PullRequestID string
	Branch        string
	AppName       string
	Image         *string
	Replicas      *int32
	Config        []*pbv1.Config
}

func (k *k8sClient) UpdateDeployment(ctx context.Context, dep *appsv1.Deployment, in UpdateDeploymentParams, opts ...Option) (*appsv1.Deployment, error) {
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	k.l.Info("updating deployment", "namespace", in.Namespace, "appName", in.AppName, "image", in.Image, CreatedByLabel, o.labelSelector[CreatedByLabel])

	if in.Image != nil {
		dep.Spec.Template.Spec.Containers[0].Image = *in.Image
	}

	if in.Replicas != nil {
		dep.Spec.Replicas = in.Replicas
	}

	if in.Config != nil {
		var config []corev1.EnvVar
		for _, c := range in.Config {
			config = append(config, corev1.EnvVar{
				Name:  c.GetName(),
				Value: c.GetValue(),
			})
		}

		dep.Spec.Template.Spec.Containers[0].Env = config
	}

	return k.client.AppsV1().Deployments(in.Namespace).Update(ctx, dep, metav1.UpdateOptions{})
}

func (k *k8sClient) DeleteDeployment(ctx context.Context, namespace string, opt ...Option) error {
	o := newOption()

	for _, opt := range opt {
		opt(o)
	}

	deps, err := k.client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.Set(o.labelSelector).String(),
	})

	if err != nil {
		return err
	}

	for _, dep := range deps.Items {
		err := k.client.AppsV1().Deployments(namespace).Delete(ctx, dep.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
