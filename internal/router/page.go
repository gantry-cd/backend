package router

import (
	"net/http"

	"github.com/gantrycd/backend/internal/handler/application/controller"
	"github.com/gantrycd/backend/internal/router/middleware"
	"github.com/gantrycd/backend/internal/usecases/application/bff"
)

func (r *router) page() {
	bc := controller.NewBff(
		bff.NewBff(r.controllerConn),
	)

	r.mux.Handle("GET /organizations",
		middleware.BuildChain(
			r.middleware.Integrate(
				http.HandlerFunc(bc.Home),
			),
		))
	r.mux.Handle("GET /organizations/{organization}/repositories", (http.HandlerFunc(bc.RepositoryApps)))
	r.mux.Handle("GET /organizations/{organization}/repositories/{repository}/pulls/{pullRequestID}", middleware.BuildChain(r.middleware.Integrate(http.HandlerFunc(bc.GetBranchInfo))))
}
