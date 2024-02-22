package router

import (
	"net/http"

	"github.com/gantrycd/backend/internal/handler/application/controller"
	"github.com/gantrycd/backend/internal/router/middleware"
	"github.com/gantrycd/backend/internal/usecases/application/resource"
)

func (r *router) Usage() {
	uc := controller.New(
		resource.New(r.controllerConn),
	)

	r.mux.Handle("GET /usage", middleware.BuildChain(
		r.middleware.Integrate(
			http.HandlerFunc(uc.Usage),
		),
	))
}
