package github

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/gantrycd/backend/cmd/config"
	ghconn "github.com/gantrycd/backend/internal/driver/github"
	"github.com/gantrycd/backend/internal/models"
	"github.com/gantrycd/backend/internal/usecases/application/controller"
	ghInteractor "github.com/gantrycd/backend/internal/usecases/application/github"
	"github.com/gantrycd/backend/internal/utils/conf"
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
	case "opened":
		ge.l.Info(fmt.Sprintf("pull request opened: %v", e.Organization.Login))
		if err := ge.pullRequestOpened(client, e); err != nil {
			ge.l.Error("error creating preview environment", "error", err.Error())
		}

	case "closed":
		ge.l.Info(fmt.Sprintf("pull request closed: %v", e))
		if err := ge.pullRequestClosed(client, e); err != nil {
			ge.l.Error("error deleting preview environment", "error", err.Error())
		}
	case "synchronize":
		ge.l.Info(fmt.Sprintf("pull request synchronize: %v", e))

	default:
		ge.l.Info(fmt.Sprintf("pull request event action not supported: %v", *e.Action))
	}
	return nil
}

func (ge *handler) pullRequestOpened(client *github.Client, e *github.PullRequestEvent) error {

	c, err := parseConfig(*e.PullRequest.Body)
	if err != nil {
		ge.l.Error("error parsing config", "error", err.Error())
		return err
	}

	ghClient, err := ghconn.GitHubConnection(
		config.Config.GitHub.AppID,
		*e.Installation.ID,
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
		Config:       *c,
		GhClient:     ghInteractor.New(ghClient),
	})
}

func (ge *handler) pullRequestClosed(client *github.Client, e *github.PullRequestEvent) error {
	return ge.interactor.DeletePreviewEnvironment(context.Background(), controller.DeletePreviewEnvironmentParams{
		Organization: *e.Organization.Login,
		Repository:   *e.Repo.Name,
		PrNumber:     fmt.Sprintf("%d", *e.Number),
		Branch:       *e.PullRequest.Head.Ref,
	})
}

func (ge *handler) pullRequestSynchronize(client *github.Client, e *github.PullRequestEvent) error {
	return ge.interactor.UpdatePreviewEnvironment(context.Background(), controller.UpdatePreviewEnvironmentParams{})
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

		log.Println(scan, line)

		if scan {
			raw += line + "\n"
		}
	}

	return conf.LoadConf(raw)
}
