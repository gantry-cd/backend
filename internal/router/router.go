package router

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gantrycd/backend/internal/router/middleware"
	v1 "github.com/gantrycd/backend/proto"
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

	r.Usage()
	return r.mux
}
