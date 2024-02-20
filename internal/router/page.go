package router

import (
	"net/http"

	"github.com/gantrycd/backend/internal/handler/application/controller"
	"github.com/gantrycd/backend/internal/router/middleware"
	"github.com/gantrycd/backend/internal/usecases/bff"
)

func (r *router) page() {
	bc := controller.NewBff(
		bff.NewBff(r.controllerConn),
	)

	r.mux.Handle("GET /organization",
		middleware.BuildChain(
			r.middleware.Integrate(
				http.HandlerFunc(bc.Home),
			),
		))
	r.mux.Handle("GET /organization/{organization}/repository", (http.HandlerFunc(bc.RepositoryApps)))
}
