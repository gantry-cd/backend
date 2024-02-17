package router

import (
	"net/http"

	"github.com/gantrycd/backend/internal/handler/application/controller"
	"github.com/gantrycd/backend/internal/router/middleware"
)

func (r *router) Usage() {
	uc := controller.New(
		r.resourceConn,
	)

	r.mux.Handle("/usage", middleware.BuildChain(http.HandlerFunc(uc.Usage)))
}
