package controller

import (
	"context"

	v1 "github.com/gantrycd/backend/proto/k8s-controller"
	"google.golang.org/protobuf/types/known/emptypb"
	appv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CustomController struct {
	v1.UnimplementedK8SCustomControllerServer
	client *kubernetes.Clientset
}

func NewCustomController(client *kubernetes.Clientset) v1.K8SCustomControllerServer {
	return &CustomController{
		client: client,
	}
}

func (c *CustomController) CreateNamespace(ctx context.Context, in *v1.CreateNamespaceRequest) (*v1.CreateNamespaceReply, error) {
	ns, err := c.client.CoreV1().Namespaces().Create(ctx, &appv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: in.Name,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return &v1.CreateNamespaceReply{
		Name: ns.Name,
	}, nil
}

func (c *CustomController) ListNamespaces(context.Context, *emptypb.Empty) (*v1.ListNamespacesReply, error) {
	panic("implement me")
}

func (c *CustomController) DeleteNamespace(context.Context, *v1.DeleteNamespaceRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (c *CustomController) ApplyDeployment(context.Context, *v1.CreateDeploymentRequest) (*v1.CreateDeploymentReply, error) {
	panic("implement me")
}

func (c *CustomController) DeleteDeployment(context.Context, *v1.DeleteDeploymentRequest) (*emptypb.Empty, error) {
	panic("implement me")
}
