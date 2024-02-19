package controller

import (
	"context"

	"github.com/gantrycd/backend/internal/models"
	v1 "github.com/gantrycd/backend/proto"
)

// CreatePreviewEnvironmentParams はプレビュー環境を作成するためのパラメータです。
type CreatePreviewEnvironmentParams struct {
	Organization string
	Repository   string
	PrNumber     string
	Branch       string

	Config models.PullRequestConfig
}

func (ge *githubAppEvents) CreatePreviewEnvironment(ctx context.Context, param CreatePreviewEnvironmentParams) (*v1.CreatePreviewReply, error) {
	// TODO: image buildする
	image := "nginx:1.16"

	var configs []*v1.Config
	for _, c := range param.Config.ConfigMaps {
		configs = append(configs, &v1.Config{
			Name:  c.Name,
			Value: c.Value,
		})
	}

	// デプロイする
	return ge.customController.CreatePreview(ctx, &v1.CreatePreviewRequest{
		Organization:  param.Organization,
		Repository:    param.Repository,
		PullRequestId: param.PrNumber,
		Branch:        param.Branch,
		Image:         image,
		Replicas:      param.Config.Replicas,
		Configs:       configs,
		ExposePorts:   param.Config.ExposePort,
	})
}

type DeletePreviewEnvironmentParams struct {
	Organization string
	Repository   string
	PrNumber     string
	Branch       string
}

func (ge *githubAppEvents) DeletePreviewEnvironment(ctx context.Context, param DeletePreviewEnvironmentParams) error {
	_, err := ge.customController.DeletePreview(ctx, &v1.DeletePreviewRequest{
		Organization:  param.Organization,
		Repository:    param.Repository,
		PullRequestId: param.PrNumber,
		Branch:        param.Branch,
	})

	return err
}
