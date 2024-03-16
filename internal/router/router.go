package router

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/aura-cd/backend/internal/router/middleware"
	v1 "github.com/aura-cd/backend/proto"
)

type router struct {
	mux        *http.ServeMux
	middleware middleware.Middleware

	l *slog.Logger

	controllerConn v1.K8SCustomControllerClient
}

func NewRouter(
	controllerConn v1.K8SCustomControllerClient,
) http.Handler {
	r := &router{
		mux:            http.NewServeMux(),
		l:              slog.New(slog.NewTextHandler(os.Stdout, nil)).WithGroup("router"),
		middleware:     middleware.NewMiddleware(),
		controllerConn: controllerConn,
	}
	r.GitHubEvent()
	r.health()
	r.page()
	r.Usage()
	r.Log()
	return r.mux
}
