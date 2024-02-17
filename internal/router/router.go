package router

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gantrycd/backend/internal/router/middleware"
	controllerV1 "github.com/gantrycd/backend/proto/k8s-controller"
	resourceV1 "github.com/gantrycd/backend/proto/metric"
)

type router struct {
	mux        *http.ServeMux
	middleware middleware.Middleware

	l *slog.Logger

	controllerConn controllerV1.K8SCustomControllerClient
	resourceConn   resourceV1.ResourceWatcherClient
}

func NewRouter(
	controllerConn controllerV1.K8SCustomControllerClient,
	resourceConn resourceV1.ResourceWatcherClient,
) http.Handler {
	r := &router{
		mux:            http.NewServeMux(),
		l:              slog.New(slog.NewTextHandler(os.Stdout, nil)).WithGroup("router"),
		controllerConn: controllerConn,
		resourceConn:   resourceConn,
	}

	return r.mux
}
