package github

import (
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v29/github"
)

func GitHubConnection(
	applicationID,
	installationID int64,
	crtPath string,
) (*github.Client, error) {
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, applicationID, installationID, crtPath)
	if err != nil {
		return nil, err
	}

	return github.NewClient(&http.Client{Transport: itr, Timeout: 5 * time.Second}), nil
}
