package controller

import (
	"context"

	"github.com/gantrycd/backend/internal/usecases/core/k8sclient"
	v1 "github.com/gantrycd/backend/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *controller) BuildImage(ctx context.Context, in *v1.BuildImageRequest) (*emptypb.Empty, error) {
	var (
		dockerfileDir  string = in.DockerfileDir
		dockerfilePath string = in.DockerfilePath
	)

	if dockerfileDir == "" {
		dockerfileDir = "."
	}

	if dockerfilePath == "" {
		dockerfilePath = "Dockerfile"
	}

	if in.ImageName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "deployment not found")
	}

	err := c.control.Builder(ctx, k8sclient.BuilderParams{
		Namespace:     in.Namespace,
		GitLink:       in.GitRepo,
		Repository:    in.Repository,
		Branch:        in.Branch,
		DockerBaseDir: dockerfileDir,
		DockrFilePath: dockerfilePath,
		ImageName:     in.ImageName,
	},
		k8sclient.WithAppLabel(k8sclient.AppLabel),
		k8sclient.WithBaseBranchLabel(in.Branch),
		k8sclient.WithPrIDLabel(in.PullRequestId),
		k8sclient.WithCreatedByLabel(k8sclient.CreatedByLabel),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
