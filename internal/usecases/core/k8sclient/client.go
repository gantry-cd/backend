package k8sclient

import (
	"context"
	"log/slog"
	"os"

	restclient "k8s.io/client-go/rest"

	protoV1 "github.com/gantrycd/backend/proto"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type k8sClient struct {
	client  *kubernetes.Clientset
	metrics *metrics.Clientset
	l       *slog.Logger
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
	CreateClientIPService(ctx context.Context, param CreateServiceNodePortParams, opts ...Option) (*corev1.Service, error)
	DeleteService(ctx context.Context, namespace string, opts ...Option) error

	// replica set
	GetReplicaSet(ctx context.Context, namespace string, prefix string) (*appsv1.ReplicaSet, error)

	// pod
	GetPods(ctx context.Context, namespace, prefix string) ([]*corev1.Pod, error)
	GetPodUsage(ctx context.Context, namespace, podName string) (*protoV1.Pod, error)

	// log
	GetLogs(namespace string, podName string, option corev1.PodLogOptions) *restclient.Request

	// builder
	Builder(ctx context.Context, param BuilderParams, opts ...Option) (string, error)
}

func New(client *kubernetes.Clientset, metrics *metrics.Clientset) K8SClient {
	return &k8sClient{
		client:  client,
		metrics: metrics,
		l:       slog.New(slog.NewTextHandler(os.Stdout, nil)).WithGroup("k8s-client"),
	}
}
