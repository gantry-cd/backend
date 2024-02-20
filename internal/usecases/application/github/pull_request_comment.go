package github

import (
	"context"

	"github.com/google/go-github/v29/github"
)

type MetaData struct {
	Organization string
	Repository   string
	Number       int
}

func (c *gitHubClientInteractor) CreateReview(ctx context.Context, meta MetaData, comment string) (*github.PullRequestReview, *github.Response, error) {
	return c.client.PullRequests.CreateReview(ctx, meta.Organization, meta.Repository, meta.Number, &github.PullRequestReviewRequest{
		Body:  &comment,
		Event: github.String("COMMENT"),
	})
}
