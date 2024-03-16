package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/aura-cd/backend/cmd/config"
	ghconn "github.com/aura-cd/backend/internal/driver/github"
	"github.com/aura-cd/backend/internal/models"
	"github.com/aura-cd/backend/internal/usecases/application/controller"
	ghInteractor "github.com/aura-cd/backend/internal/usecases/application/github"
	"github.com/aura-cd/backend/internal/utils/conf"
	"github.com/google/go-github/v29/github"
)

func (ge *handler) PullRequest(e *github.PullRequestEvent) error {
	client, err := ghconn.GitHubConnection(
		config.Config.GitHub.AppID,
		*e.Installation.ID,
		config.Config.GitHub.CrtPath,
	)
	if err != nil {
		ge.l.Error("error creating github client", "error", err.Error())
		return err
	}

	switch *e.Action {
	case PullRequestOpened:
		ge.l.Info(fmt.Sprintf("pull request opened: %v", e.Organization.Login))
		if err := ge.pullRequestOpened(client, e); err != nil {
			ge.l.Error("error creating preview environment", "error", err.Error())
		}
	case PullRequestClosed:
		ge.l.Info(fmt.Sprintf("pull request closed: %v", e))
		if err := ge.pullRequestClosed(client, e); err != nil {
			ge.l.Error("error deleting preview environment", "error", err.Error())
		}
	case PullRequestSynchronize:
		ge.l.Info(fmt.Sprintf("pull request synchronize: %v", e))
		if err := ge.pullRequestSynchronize(client, e); err != nil {
			ge.l.Error("error updating preview environment", "error", err.Error())
		}
	default:
		ge.l.Info(fmt.Sprintf("pull request event action not supported: %v", *e.Action))
	}
	return nil
}

// pullRequestOpened はプルリクエストがオープンされたときにプレビュー環境を作成する
func (ge *handler) pullRequestOpened(client *github.Client, e *github.PullRequestEvent) error {
	// Body から設定を取得して設定をパースする
	c, err := parseConfig(e.PullRequest.GetBody())
	if err != nil {
		ge.l.Error("error parsing config", "error", err.Error())
		return err
	}

	installID := e.Installation.GetID()

	// GitHub クライアントを作成する
	ghClient, err := ghconn.GitHubConnection(
		config.Config.GitHub.AppID,
		installID,
		config.Config.GitHub.CrtPath,
	)
	if err != nil {
		ge.l.Error("error parsing config", "error", err.Error())
		return err
	}

	return ge.interactor.CreatePreviewEnvironment(context.Background(), controller.CreatePreviewEnvironmentParams{
		Organization: *e.Organization.Login,
		Repository:   *e.Repo.Name,
		PrNumber:     *e.Number,
		Branch:       *e.PullRequest.Head.Ref,
		GitLink:      *e.PullRequest.Head.Repo.CloneURL,
		GhInstallID:  installID,
		Config:       *c,
		GhClient:     ghInteractor.New(ghClient),
	})
}

// pullRequestClosed はプルリクエストがクローズされたときにプレビュー環境を削除する
func (ge *handler) pullRequestClosed(client *github.Client, e *github.PullRequestEvent) error {
	return ge.interactor.DeletePreviewEnvironment(context.Background(), controller.DeletePreviewEnvironmentParams{
		Organization: *e.Organization.Login,
		Repository:   *e.Repo.Name,
		PrNumber:     fmt.Sprintf("%d", *e.Number),
		Branch:       *e.PullRequest.Head.Ref,
	})
}

// pullRequestSynchronize はプルリクエストが更新されたときにプレビュー環境を更新する
func (ge *handler) pullRequestSynchronize(client *github.Client, e *github.PullRequestEvent) error {
	c, err := parseConfig(e.PullRequest.GetBody())
	if err != nil {
		ge.l.Error("error parsing config", "error", err.Error())
		return err
	}

	installID := e.Installation.GetID()

	// GitHub クライアントを作成する
	ghClient, err := ghconn.GitHubConnection(
		config.Config.GitHub.AppID,
		installID,
		config.Config.GitHub.CrtPath,
	)
	if err != nil {
		ge.l.Error("error parsing config", "error", err.Error())
		return err
	}

	return ge.interactor.UpdatePreviewEnvironment(context.Background(), controller.UpdatePreviewEnvironmentParams{
		Organization: *e.Organization.Login,
		Repository:   *e.Repo.Name,
		PrNumber:     *e.Number,
		Branch:       *e.PullRequest.Head.Ref,
		GitLink:      *e.PullRequest.Head.Repo.CloneURL,
		GhInstallID:  installID,
		Config:       *c,
		GhClient:     ghInteractor.New(ghClient),
	})
}

const (
	configureIndentifer = "<%sgantry-config>"
)

// parseConfig はプルリクエストのメッセージから設定をパースする
func parseConfig(prMessage string) (*models.PullRequestConfig, error) {
	var (
		scan bool
		raw  string
	)

	for _, line := range strings.Split(prMessage, "\r\n") {
		if strings.HasPrefix(line, fmt.Sprintf(configureIndentifer, "")) {
			scan = true
			continue
		}
		if strings.HasPrefix(line, fmt.Sprintf(configureIndentifer, "/")) {
			scan = false
			continue
		}

		if scan {
			raw += strings.TrimSpace(line) + "\n"
		}
	}

	return conf.LoadConf(raw)
}
