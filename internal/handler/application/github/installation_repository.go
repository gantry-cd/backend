package github

import (
	"fmt"

	"github.com/google/go-github/v29/github"
)

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
