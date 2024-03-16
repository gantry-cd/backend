package main

import (
	"fmt"
	"os"

	"github.com/aura-cd/backend/cmd/config"
	"github.com/aura-cd/backend/internal/driver/pbclient"
	"github.com/aura-cd/backend/internal/router"
	"github.com/aura-cd/backend/internal/server/http"
	v1 "github.com/aura-cd/backend/proto"
)

func init() {
	config.LoadEnv(
		".env/keycloak.env",
		".env/bff.env",
		".env/github.env",
		".env/harbor.env",
	)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

func run() error {
	controllerPbc := pbclient.NewConn(config.Config.Bff.K8SControllerAddr)

	if err := controllerPbc.Connect(); err != nil {
		return fmt.Errorf("failed to connect to k8s controller: %w", err)
	}
	defer controllerPbc.Close()

	handler := router.NewRouter(
		v1.NewK8SCustomControllerClient(controllerPbc.Client()),
	)

	server := http.New(
		handler,
		http.WithPort(fmt.Sprintf("%d", config.Config.Bff.Port)),
		http.WithHost(config.Config.Bff.Host),
		http.WithShutdownTimeout(10),
	)

	return server.Run()
}
