package controller

import (
	"context"

	"github.com/gantrycd/backend/internal/usecases/core/k8sclient"
	v1 "github.com/gantrycd/backend/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *controller) BuildImage(ctx context.Context, in *v1.BuildImageRequest) (*v1.BuildImageReply, error) {
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

	image, err := c.control.Builder(context.Background(), k8sclient.BuilderParams{
		Namespace: in.Namespace,
		Branch:    in.Branch,
		BuilderParam: k8sclient.ImageBuilderParams{
			Repository:     in.Repository,
			GitRepo:        in.GitRepo,
			GitBranch:      in.Branch,
			ImageName:      in.ImageName,
			DockerBaseDir:  dockerfileDir,
			DockerFilePath: dockerfilePath,
		},
	},
		k8sclient.WithAppLabel(k8sclient.AppLabel),
		k8sclient.WithBaseBranchLabel(in.Branch),
		k8sclient.WithPrIDLabel(in.PullRequestId),
		k8sclient.WithCreatedByLabel(k8sclient.CreatedByLabel),
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &v1.BuildImageReply{
		Image: image,
	}, nil
}
