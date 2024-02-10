package controller

import (
	"context"

	v1 "github.com/gantrycd/backend/proto/k8s-controller"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CustomController struct {
	v1.UnimplementedK8SCustomControllerServer
}

func NewCustomController() v1.K8SCustomControllerServer {
	return &CustomController{}
}

func (c *CustomController) CreateNamespace(context.Context, *v1.CreateNamespaceRequest) (*v1.CreateNamespaceReply, error) {
	panic("implement me")
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
