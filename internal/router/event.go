package router

import (
	"net/http"

	"github.com/aura-cd/backend/internal/adapters/controllers/application/github"
	"github.com/aura-cd/backend/internal/usecases/application/controller"
)

func (r *router) GitHubEvent() {
	githubHandler := github.New(
		controller.New(r.controllerConn),
	)
	r.mux.Handle("POST /github/app/webhook", http.HandlerFunc(githubHandler.GithubAppsHandler))
}
