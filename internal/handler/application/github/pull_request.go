package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/gantrycd/backend/cmd/config"
	ghconn "github.com/gantrycd/backend/internal/driver/github"
	"github.com/gantrycd/backend/internal/models"
	"github.com/gantrycd/backend/internal/usecases/application/githubapp"
	"github.com/gantrycd/backend/internal/utils/conf"
	"github.com/google/go-github/v29/github"
)

func (ge *handler) PullRequest(e *github.PullRequestEvent) error {
	ge.l.Info(fmt.Sprintf("pull request event received: %v", e))
	ctx := context.Background()

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
		if err := ge.pullRequestOpend(client, e); err != nil {
			ge.l.Error("error creating preview environment", "error", err.Error())
		}

	case "closed":
		ge.l.Info(fmt.Sprintf("pull request closed: %v", e))
		if err := ge.interactor.DeletePreviewEnvironment(ctx, githubapp.DeletePreviewEnvironmentParams{
			Organization: *e.Organization.Login,
			Repository:   *e.Repo.Name,
			PrNumber:     fmt.Sprintf("%d", *e.Number),
			Branch:       *e.PullRequest.Head.Ref,
		}); err != nil {
			ge.l.Error("error deleting preview environment", "error", err.Error())
		}
	default:
		ge.l.Info(fmt.Sprintf("pull request event action not supported: %v", *e.Action))
	}
	return nil
}

func (ge *handler) pullRequestOpend(client *github.Client, e *github.PullRequestEvent) error {
	config, err := parseConfig(*e.PullRequest.Body)
	if err != nil {
		ge.l.Error("error parsing config", "error", err.Error())
		return err
	}
	dep, err := ge.interactor.CreatePreviewEnvironment(context.Background(), githubapp.CreatePreviewEnvironmentPrarams{
		Organization: *e.Organization.Login,
		Repository:   *e.Repo.Name,
		PrNumber:     fmt.Sprintf("%d", *e.Number),
		Branch:       *e.PullRequest.Head.Ref,
		Config:       *config,
	})
	if err != nil {
		return fmt.Errorf("error creating preview environment: %w", err)
	}

	var message string
	for _, np := range dep.NodePorts {
		message += fmt.Sprintf("port %d: %d\n", np.Port, np.Target)
	}

	prReview, ghResp, err := client.PullRequests.CreateReview(context.Background(), *e.Organization.Login, *e.Repo.Name, *e.Number, &github.PullRequestReviewRequest{
		Body:  github.String("deploy preview environment create successful! 🚀 \n " + message),
		Event: github.String("APPROVE"),
	})
	if err != nil {
		ge.l.Error("error creating review", "error", err.Error())
		return err
	}

	ge.l.Info(fmt.Sprintf("pull request review created: %v", prReview))
	ge.l.Info(fmt.Sprintf("github response: %v", ghResp))

	return nil
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

	for _, line := range strings.Split(prMessage, "\n") {
		if strings.HasPrefix(line, configureIndentifer) || strings.HasPrefix(line, fmt.Sprintf(configureIndentifer, "/")) {
			scan = true
			continue
		}

		if scan {
			raw += line + "\n"
		}
	}

	return conf.LoadConf(raw)
}
