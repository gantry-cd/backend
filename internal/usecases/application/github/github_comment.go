package github

import (
	"context"

	"github.com/google/go-github/v29/github"
)

type gitHubClientInteractor struct {
	client *github.Client
}

type GitHubClientInteractor interface {
	CreatePullRequestComment(ctx context.Context, meta MetaData, comment string) (*github.PullRequestComment, *github.Response, error)
	EditPullRequestComment(ctx context.Context, meta MetaData, param EditPullRequestComment) (*github.PullRequestComment, *github.Response, error)
}

func New(client *github.Client) GitHubClientInteractor {
	return &gitHubClientInteractor{client: client}
}
