package router

import (
	"net/http"

	"github.com/aura-cd/backend/internal/handler/application/log"
	"github.com/aura-cd/backend/internal/router/middleware"
	"github.com/aura-cd/backend/internal/usecases/application/logger"
)

func (r *router) Log() {
	uc := log.New(
		logger.New(r.controllerConn),
	)

	r.mux.Handle("/organizations/{organization}/pods/{pod}/log", middleware.BuildChain(
		r.middleware.Integrate(
			http.HandlerFunc(uc.Log),
		),
	))
}
