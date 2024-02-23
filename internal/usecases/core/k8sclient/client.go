package k8sclient

import (
	"context"
	"log/slog"
	"os"

	restclient "k8s.io/client-go/rest"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	CreateDeployment(ctx context.Context, in CreateDeploymentParams, opts ...Option) (*appsv1.Deployment, error)
	GetDeployment(ctx context.Context, param GetDeploymentParams) (*appsv1.Deployment, error)
	ListDeployments(ctx context.Context, namespace string, opts ...Option) (*appsv1.DeploymentList, error)
	UpdateDeployment(ctx context.Context, dep *appsv1.Deployment, in UpdateDeploymentParams, opts ...Option) (*appsv1.Deployment, error)
	DeleteDeployment(ctx context.Context, namespace string, opts ...Option) error

	// service
	CreateLoadBalancerService(ctx context.Context, param CreateServiceNodePortParams, opts ...Option) (*corev1.Service, error)
	DeleteService(ctx context.Context, namespace string, opts ...Option) error

	// replica set
	GetReplicaSet(ctx context.Context, namespace string, prefix string) (*appsv1.ReplicaSet, error)

	// pod
	GetPods(ctx context.Context, namespace, prefix string) ([]*corev1.Pod, error)
	CreatePod(ctx context.Context, namespace string, pod *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error)
	UpdatePod(ctx context.Context, namespace string, pod *corev1.Pod, opts metav1.UpdateOptions) (*corev1.Pod, error)
	DeletePod(ctx context.Context, namespace, name string, opts metav1.DeleteOptions) error

	// log
	GetLogs(namespace string, podName string, option corev1.PodLogOptions) *restclient.Request

	// builder
	Builder(ctx context.Context, param BuilderParams, opts ...Option) (string, error)

	// config map
	CreateConfigMap(ctx context.Context, namespace string, name string, data map[string]string, opts metav1.CreateOptions) error
	UpdateConfigMap(ctx context.Context, namespace string, name string, data map[string]string, opts metav1.UpdateOptions) error
	DeleteConfigMap(ctx context.Context, namespace string, name string, opts metav1.DeleteOptions) error
}

func New(client *kubernetes.Clientset) K8SClient {
	return &k8sClient{
		client: client,
		l:      slog.New(slog.NewTextHandler(os.Stdout, nil)).WithGroup("k8s-client"),
	}
}
