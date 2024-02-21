package controller

import (
	"context"
	"errors"

	coreErr "github.com/gantrycd/backend/internal/error"
	"github.com/gantrycd/backend/internal/usecases/core/k8sclient"
	"github.com/gantrycd/backend/internal/utils"
	"github.com/gantrycd/backend/internal/utils/branch"
	v1 "github.com/gantrycd/backend/proto"
	"google.golang.org/protobuf/types/known/emptypb"
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

	service, err := c.control.CreateNodePortService(ctx,
		k8sclient.CreateServiceNodePortParams{
			Namespace:   in.Organization,
			ServiceName: in.Repository,
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

	var nodePorts []*v1.NodePort
	for _, port := range service.Spec.Ports {
		nodePorts = append(nodePorts, &v1.NodePort{
			Port:   port.Port,
			Target: port.NodePort,
		})
	}

	return &v1.CreatePreviewReply{
		Name:      deps.Name,
		Namespace: deps.Namespace,
		Version:   deps.ResourceVersion,
		NodePorts: nodePorts,
	}, nil
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

	return new(emptypb.Empty), nil
}
