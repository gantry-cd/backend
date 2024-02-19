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

func (c *gitHubClientInteractor) CreatePullRequestComment(ctx context.Context, meta MetaData, comment string) (*github.PullRequestComment, *github.Response, error) {
	return c.client.PullRequests.CreateComment(ctx, meta.Organization, meta.Repository, meta.Number, &github.PullRequestComment{
		Body: &comment,
	})
}

type EditPullRequestComment struct {
	CommentID int64
	Comment   string
}

func (c *gitHubClientInteractor) EditPullRequestComment(ctx context.Context, meta MetaData, param EditPullRequestComment) (*github.PullRequestComment, *github.Response, error) {
	return c.client.PullRequests.EditComment(ctx, meta.Organization, meta.Repository, param.CommentID, &github.PullRequestComment{
		Body: &param.Comment,
	})
}
