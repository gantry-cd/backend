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
	resourceConn   v1.ResourceWatcherClient
}

func NewRouter(
	controllerConn v1.K8SCustomControllerClient,
	resourceConn v1.ResourceWatcherClient,
) http.Handler {
	r := &router{
		mux:            http.NewServeMux(),
		l:              slog.New(slog.NewTextHandler(os.Stdout, nil)).WithGroup("router"),
		controllerConn: controllerConn,
		resourceConn:   resourceConn,
	}

	return r.mux
}
