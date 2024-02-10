package webhook

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v29/github"
)

const (
	organization = iota
	repository
)

func (ge *handler) Installation(e *github.InstallationEvent) error {
	ge.l.Info(fmt.Sprintf("installation event received: %v", e))
	ctx := context.Background()

	for _, repo := range e.Repositories {
		ge.l.Info(fmt.Sprintf("repository: %v", repo))

		names := strings.Split(*repo.FullName, "/")

		nss, err := ge.interactor.ListNameSpace(ctx, names[organization])
		if err != nil {
			ge.l.Error("error listing namespaces", "error", err.Error())
			return err
		}

		if !isInclude(nss, names[repository]) {
			if err := ge.interactor.CreateNameSpace(ctx, names[organization]); err != nil {
				ge.l.Error("error creating namespace", "error", err.Error())
				return err
			}
		}
	}

	return nil
}

func (ge *handler) InstallationRepositories(e *github.InstallationRepositoriesEvent) error {
	ge.l.Info(fmt.Sprintf("installation repositories event received: %v", e))
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
