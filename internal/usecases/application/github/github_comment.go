package github

import (
	"context"

	"github.com/google/go-github/v29/github"
)

type gitHubClientInteractor struct {
	client *github.Client
}

type GitHubClientInteractor interface {
	CreateReview(ctx context.Context, meta MetaData, comment string) (*github.PullRequestReview, *github.Response, error)
	GetToken(ctx context.Context, installID int64) (string, error)
}

func New(client *github.Client) GitHubClientInteractor {
	return &gitHubClientInteractor{client: client}
}
