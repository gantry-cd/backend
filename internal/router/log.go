package router

import (
	"net/http"

	"github.com/gantrycd/backend/internal/handler/application/log"
	"github.com/gantrycd/backend/internal/router/middleware"
	"github.com/gantrycd/backend/internal/usecases/application/logger"
)

func (r *router) Log() {
	uc := log.New(
		logger.New(r.controllerConn),
	)

	r.mux.Handle("/organizations/{organization}/repositories/{repository}/pulls/{pull}/pods/{pod}/log", middleware.BuildChain(
		r.middleware.Integrate(
			http.HandlerFunc(uc.Log),
		),
	))
}
