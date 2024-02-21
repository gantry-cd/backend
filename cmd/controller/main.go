package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"

	"github.com/gantrycd/backend/cmd/config"
	"github.com/gantrycd/backend/internal/driver/k8s"
	"github.com/gantrycd/backend/internal/handler/core/controller"
	"github.com/gantrycd/backend/internal/usecases/core/resource"
	v1 "github.com/gantrycd/backend/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	config.LoadEnv(
		".env/controller.env",
		".env/harbor.env",
		".env/github.env",
	)
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}

func run() error {
	path := ".kube/config"
	k8sConn := k8s.New(&path)
	client, err := k8sConn.ConnectTyped("")
	if err != nil {
		return err
	}

	metric, err := k8sConn.ConnectMetrics("")
	if err != nil {
		return err
	}

	l := slog.New(slog.NewTextHandler(os.Stderr, nil)).WithGroup("server")

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Config.Controller.Host, config.Config.Controller.Port))
	if err != nil {
		return err
	}

	// TODO: Implement the server
	v1.RegisterK8SCustomControllerServer(grpcServer, controller.NewController(client, resource.New(metric)))

	reflection.Register(grpcServer)

	go func() {
		l.Info(fmt.Sprintf("server is running at %s", fmt.Sprintf("%s:%d", config.Config.Controller.Host, config.Config.Controller.Port)))
		if err := grpcServer.Serve(listener); err != nil {
			l.Error("failed to serve", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	l.Info("stopping gRPC server...")
	grpcServer.GracefulStop()

	return nil
}
