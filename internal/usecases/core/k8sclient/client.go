package k8sclient

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

type k8sClient struct {
	client *kubernetes.Clientset

	l *slog.Logger
}

type K8SClient interface {
	// namespace
	CreateNamespace(ctx context.Context, name string, opts ...Option) (*corev1.Namespace, error)
	ListNamespaces(ctx context.Context, opts ...Option) (*corev1.NamespaceList, error)
	DeleteNamespace(ctx context.Context, name string) error

	// deployment
	CreateDeployment(ctx context.Context, namespace, podName, image string, opts ...Option) (*appsv1.Deployment, error)
	GetDeployment(ctx context.Context, namespace, repository, prID string) (*appsv1.Deployment, error)
	ListDeployments(ctx context.Context, namespace string, opts ...Option) (*appsv1.DeploymentList, error)
	DeleteDeployment(ctx context.Context, namespace, repository, prID string) error

	//
}

func New(client *kubernetes.Clientset) K8SClient {
	return &k8sClient{
		client: client,
		l:      slog.New(slog.NewTextHandler(os.Stdout, nil)).WithGroup("k8s-client"),
	}
}

func (k *k8sClient) CreateNamespace(ctx context.Context, name string, opts ...Option) (*corev1.Namespace, error) {
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	k.l.Info("creating namespace", "name", name, CreatedByLabel, o.labelSelector[CreatedByLabel])

	return k.client.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: o.labelSelector,
		},
	}, metav1.CreateOptions{})
}

func (k *k8sClient) ListNamespaces(ctx context.Context, opts ...Option) (*corev1.NamespaceList, error) {
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	return k.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{LabelSelector: labels.Set(o.labelSelector).String()})
}

func (k *k8sClient) DeleteNamespace(ctx context.Context, name string) error {
	return k.client.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
}

func (k *k8sClient) CreateDeployment(ctx context.Context, namespace, podName, image string, opts ...Option) (*appsv1.Deployment, error) {
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	o.labelSelector[AppLabel] = podName

	return k.client.AppsV1().Deployments(namespace).Create(ctx, &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: o.replica,
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
							Name:            podName,
							Image:           image,
							ImagePullPolicy: o.containerOption[podName].imagePullPolicy,
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{})
}

func (k *k8sClient) GetDeployment(ctx context.Context, namespace, repository, prID string) (*appsv1.Deployment, error) {
	deps, err := k.client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("repository=%s,pr-number=%s", repository, prID),
	})

	if err != nil {
		return nil, err
	}

	for _, dep := range deps.Items {
		if dep.ObjectMeta.Labels[RepositryLabel] == repository && dep.ObjectMeta.Labels[PrIDLabel] == prID {
			return &dep, nil
		}
	}

	return nil, fmt.Errorf("deployment not found")
}

func (k *k8sClient) ListDeployments(ctx context.Context, namespace string, opts ...Option) (*appsv1.DeploymentList, error) {
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	return k.client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{LabelSelector: labels.Set(o.labelSelector).String()})
}

func (k *k8sClient) DeleteDeployment(ctx context.Context, namespace, repository, prID string) error {
	deps, err := k.client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("repository=%s,pr-number=%s", repository, prID),
	})

	if err != nil {
		return err
	}

	for _, dep := range deps.Items {
		// delete deployment
		err := k.client.AppsV1().Deployments(namespace).Delete(ctx, dep.Name, metav1.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}
