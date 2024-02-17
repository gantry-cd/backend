package resource

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gantrycd/backend/internal/models"
	v1 "github.com/gantrycd/backend/proto/metric"
	"google.golang.org/protobuf/types/known/emptypb"
)

type resrouceInteractor struct {
	resource v1.ResourceWatcherClient
}

type ResrouceInteractor interface {
}

func New(resource v1.ResourceWatcherClient) ResrouceInteractor {
	return &resrouceInteractor{
		resource: resource,
	}
}

func (r *resrouceInteractor) GetResourceSSE(ctx context.Context, w http.ResponseWriter, request models.UsageRequest) error {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, _ := w.(http.Flusher)

	go func() {
		for {
			var resp models.UsageResponse
			var usages []models.Usage

			result, err := r.resource.GetResource(ctx, &emptypb.Empty{})
			if err != nil {
				continue
			}
			// Prometheusとかない場合
			if result.IsDisable {
				resp.IsDisable = true
			}

			resources := result.GetResources()
			for _, resource := range resources {

			}

			if err := json.NewEncoder(w).Encode(resp); err != nil {
				continue
			}

			flusher.Flush()
		}
	}()

	return nil
}
