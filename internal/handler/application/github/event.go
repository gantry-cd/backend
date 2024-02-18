package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/gantrycd/backend/internal/usecases/application/githubapp"
	"github.com/google/go-github/v29/github"
)

const (
	organization = iota
	repository
)

func (ge *handler) Installation(e *github.InstallationEvent) error {
	ge.l.Info(fmt.Sprintf("installation event received: %v", e))
	ctx := context.Background()

	var orgs []string
	// get all organizations
	for _, repo := range e.Repositories {
		ge.l.Info(fmt.Sprintf("repository: %v", repo))

		names := strings.Split(*repo.FullName, "/")

		if isInclude(orgs, names[organization]) {
			continue
		}

		orgs = append(orgs, names[organization])
	}

	nss, err := ge.interactor.ListNameSpace(ctx, "")
	if err != nil {
		ge.l.Error("error listing namespaces", "error", err.Error())
		return err
	}

	switch *e.Action {
	case InstallationCreated:
		for _, org := range orgs {
			orgname := strings.ToLower(org)
			if !isInclude(nss, orgname) {
				if err := ge.interactor.CreateNameSpace(ctx, org); err != nil {
					ge.l.Error("error creating namespace", "error", err.Error())
				}
			}
		}
	case InstallationDeleted:
		for _, org := range orgs {
			orgname := strings.ToLower(org)
			if isInclude(nss, orgname) {
				if err := ge.interactor.DeleteNameSpace(ctx, orgname); err != nil {
					ge.l.Error("error deleting namespace", "error", err.Error())
				}
			}
		}
	default:
		ge.l.Info(fmt.Sprintf("installation event action not supported: %v", *e.Action))
	}

	return nil
}

func (ge *handler) InstallationRepositories(e *github.InstallationRepositoriesEvent) error {
	ge.l.Info(fmt.Sprintf("installation repositories event received: %v", e))

	// ctx := context.Background()

	// リポジトリから削除された場合そのラベルを持つDeployたちを削除する
	if *e.Action == "removed" {
		for _, repo := range e.RepositoriesRemoved {
			ge.l.Info(fmt.Sprintf("repository removed: %v", repo))
		}
	}
	return nil
}

func (ge *handler) Meta(e *github.MetaEvent) error {
	ge.l.Info(fmt.Sprintf("meta event received: %v", e))
	return nil
}

func (ge *handler) Create(e *github.CreateEvent) error {
	ge.l.Info(fmt.Sprintf("create event received: %v", e))
	return nil
}

func (ge *handler) Delete(e *github.DeleteEvent) error {
	ge.l.Info(fmt.Sprintf("delete event received: %v", e))
	return nil
}

func (ge *handler) Push(e *github.PushEvent) error {
	ge.l.Info(fmt.Sprintf("push event received: %v", e))
	return nil
}

func (ge *handler) PullRequest(e *github.PullRequestEvent) error {
	ge.l.Info(fmt.Sprintf("pull request event received: %v", e))
	ctx := context.Background()

	switch *e.Action {
	case "opened":
		ge.l.Info(fmt.Sprintf("pull request opened: %v", e.Organization.Login))

		if err := ge.interactor.CreatePreviewEnvironment(ctx, githubapp.CreatePreviewEnvironmentPrarams{
			Organization: *e.Organization.Login,
			Repository:   *e.Repo.Name,
			PrNumber:     fmt.Sprintf("%d", *e.Number),
			Branch:       *e.PullRequest.Head.Ref,
		}); err != nil {
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

func (ge *handler) Repository(e *github.RepositoryEvent) error {
	ge.l.Info(fmt.Sprintf("repository event received: %v", e))
	return nil
}

func isInclude(ns []string, name string) bool {
	for _, n := range ns {
		if n == name {
			return true
		}
	}
	return false
}
