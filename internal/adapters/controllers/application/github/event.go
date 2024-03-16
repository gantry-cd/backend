package github

import (
	"fmt"

	"github.com/google/go-github/v29/github"
)

const (
	organization = iota
	repository
)

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
