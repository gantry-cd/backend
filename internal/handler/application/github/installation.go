package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v29/github"
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
