package main

import (
	"fmt"
	"os"

	"github.com/gantrycd/backend/internal/driver/pbclient"
	"github.com/gantrycd/backend/internal/router"
	"github.com/gantrycd/backend/internal/server/http"
	v1 "github.com/gantrycd/backend/proto"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

func run() error {
	controllerPbc := pbclient.NewConn(os.Getenv("K8S_CONTROLLER_ADDR"))

	if err := controllerPbc.Connect(); err != nil {
		return fmt.Errorf("failed to connect to k8s controller: %w", err)
	}
	defer controllerPbc.Close()

	handler := router.NewRouter(
		v1.NewK8SCustomControllerClient(controllerPbc.Client()),
	)

	server := http.New(
		handler,
		http.WithPort(os.Getenv("PORT")),
		http.WithHost(os.Getenv("HOST")),
		http.WithShutdownTimeout(10),
	)

	return server.Run()
}
