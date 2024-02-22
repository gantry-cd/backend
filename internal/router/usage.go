package router

import (
	"net/http"

	controller "github.com/gantrycd/backend/internal/handler/application/usage"
	"github.com/gantrycd/backend/internal/router/middleware"
	usecase "github.com/gantrycd/backend/internal/usecases/application/usage"
)

func (r *router) Usage() {
	uc := controller.New(
		usecase.New(r.controllerConn),
	)

	r.mux.Handle("GET /organizations/pods", middleware.BuildChain(
		r.middleware.Integrate(
			http.HandlerFunc(uc.Usage),
		),
	))
}
