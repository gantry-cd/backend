package main

import (
	"fmt"
	"os"

	"github.com/gantrycd/backend/internal/handler/webhook"
	"github.com/gantrycd/backend/internal/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

func run() error {
	handler := webhook.New()
	server := server.New(
		handler,
		server.WithPort(os.Getenv("PORT")),
		server.WithHost(os.Getenv("HOST")),
		server.WithShutdownTimeout(10),
	)
	return server.Run()
}
