package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/gantrycd/backend/internal/models"
	"github.com/gantrycd/backend/internal/usecases/application/github"
	v1 "github.com/gantrycd/backend/proto"
)

// CreatePreviewEnvironmentParams はプレビュー環境を作成するためのパラメータです。
type CreatePreviewEnvironmentParams struct {
	Organization string
	Repository   string
	PrNumber     int
	Branch       string

	GhClient github.GitHubClientInteractor
	Config   models.PullRequestConfig
}

func (ge *githubAppEvents) CreatePreviewEnvironment(ctx context.Context, param CreatePreviewEnvironmentParams) error {
	meta := github.MetaData{
		Organization: param.Organization,
		Repository:   param.Repository,
		Number:       param.PrNumber,
	}
	comment, _, err := param.GhClient.CreatePullRequestComment(ctx, meta, fmt.Sprintf("[%v]Creating preview environment for %s 🚀 ...\n Building Docker image at %s", time.Now().Format(time.DateTime), param.Branch, param.Config.BuildFilePath))
	if err != nil {
		return err
	}

	// TODO: image buildする
	image := "nginx:1.16"

	comment, _, err = param.GhClient.EditPullRequestComment(ctx, meta, github.EditPullRequestComment{
		CommentID: comment.GetID(),
		Comment:   fmt.Sprintf("%s\n[%v] ✔️Docker image built at %s\n Creating deployment and service... ", comment.GetBody(), time.Now().Format(time.DateTime), image),
	})
	if err != nil {
		_, _, err := param.GhClient.EditPullRequestComment(ctx, meta, github.EditPullRequestComment{
			CommentID: comment.GetID(),
			Comment:   fmt.Sprintf("%s\n[%v] ❌Failed to build Docker image: %v", comment.GetBody(), time.Now().Format(time.DateTime), err),
		})
		return err
	}

	var configs []*v1.Config
	for _, c := range param.Config.ConfigMaps {
		configs = append(configs, &v1.Config{
			Name:  c.Name,
			Value: c.Value,
		})
	}

	if param.Config.Replicas == 0 {
		param.Config.Replicas = 1
	}

	// デプロイする
	dep, err := ge.customController.CreatePreview(ctx, &v1.CreatePreviewRequest{
		Organization:  param.Organization,
		Repository:    param.Repository,
		PullRequestId: fmt.Sprintf("%d", param.PrNumber),
		Branch:        param.Branch,
		Image:         image,
		Replicas:      param.Config.Replicas,
		Configs:       configs,
		ExposePorts:   param.Config.ExposePort,
	})

	if err != nil {
		_, _, err = param.GhClient.EditPullRequestComment(ctx, meta, github.EditPullRequestComment{
			CommentID: comment.GetID(),
			Comment:   fmt.Sprintf("%s\n[%v] ❌Failed to create deployment and service: %v", comment.GetBody(), time.Now().Format(time.DateTime), err),
		})
		return err
	}

	var ports string
	for _, p := range dep.NodePorts {
		ports += fmt.Sprintf("%d:%d ", p.Port, p.Target)
	}

	_, _, err = param.GhClient.EditPullRequestComment(ctx, meta, github.EditPullRequestComment{
		CommentID: comment.GetID(),
		Comment:   fmt.Sprintf("%s\n[%v] ✔️Deployment and service created\n Exposing service on port %s ", comment.GetBody(), time.Now().Format(time.DateTime), ports),
	})
	if err != nil {
		return err
	}

	return nil
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
