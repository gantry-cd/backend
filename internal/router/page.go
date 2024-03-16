package router

import (
	"net/http"

	"github.com/aura-cd/backend/internal/adapters/controllers/application/api"
	"github.com/aura-cd/backend/internal/router/middleware"
	"github.com/aura-cd/backend/internal/usecases/bff"
)

func (r *router) page() {
	bc := api.New(
		bff.NewBff(r.controllerConn),
	)

	r.mux.Handle("GET /organizations",
		middleware.BuildChain(
			r.middleware.Integrate(
				http.HandlerFunc(bc.Home),
			),
		))
	r.mux.Handle("GET /organizations/{organization}/repositories", (http.HandlerFunc(bc.RepositoryApps)))
	r.mux.Handle("GET /organizations/{organization}/repositories/{repository}/branches", (http.HandlerFunc(bc.GetRepoBranches)))
}
