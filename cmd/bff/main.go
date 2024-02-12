package main

import (
	"fmt"
	"os"

	"github.com/gantrycd/backend/internal/driver/pbclient"
	"github.com/gantrycd/backend/internal/router"
	"github.com/gantrycd/backend/internal/server/http"
	v1 "github.com/gantrycd/backend/proto/k8s-controller"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

func run() error {
	pbc := pbclient.NewConn(os.Getenv("K8S_CONTROLLER_ADDR"))

	if err := pbc.Connect(); err != nil {
		return fmt.Errorf("failed to connect to k8s controller: %w", err)
	}
	defer pbc.Close()

	handler := router.NewRouter(
		v1.NewK8SCustomControllerClient(pbc.Client()),
	)

	server := http.New(
		handler,
		http.WithPort(os.Getenv("PORT")),
		http.WithHost(os.Getenv("HOST")),
		http.WithShutdownTimeout(10),
	)

	return server.Run()
}
