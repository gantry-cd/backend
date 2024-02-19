package router

import (
	"net/http"

	"github.com/gantrycd/backend/internal/handler/application/github"
	"github.com/gantrycd/backend/internal/usecases/application/githubapp"
)

func (r *router) GitHubEvent() {
	githubHandler := github.New(
		githubapp.New(r.controllerConn),
	)
	r.mux.Handle("POST /github/app/webhook", http.HandlerFunc(githubHandler.GithubAppsHandler))
}
