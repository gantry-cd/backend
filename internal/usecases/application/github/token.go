package github

import (
	"context"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/gantrycd/backend/cmd/config"
)

func (c *gitHubClientInteractor) GetToken(ctx context.Context, installID int64) (string, error) {
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, config.Config.GitHub.AppID, installID, config.Config.GitHub.CrtPath)
	if err != nil {
		return "", err
	}

	return itr.Token(ctx)
}
