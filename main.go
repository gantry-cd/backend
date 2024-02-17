package main

import (
	"context"
	"strings"
	"time"

	"github.com/gantrycd/backend/internal/driver/k8s"
	"github.com/gantrycd/backend/internal/usecases/core/resource"
)

func main() {
	config := "./.kube/config"
	client := k8s.New(&config)
	metr, err := client.ConnectMetrics("")
	if err != nil {
		panic(err)
	}
	r := resource.New(metr)

	for {
		println(strings.Repeat("=", 50))
		if err := r.GetLoads(context.Background(), "test-space"); err != nil {
			panic(err)
		}

		time.Sleep(5 * time.Second)
	}
}
