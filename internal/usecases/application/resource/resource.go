package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gantrycd/backend/internal/models"
	v1 "github.com/gantrycd/backend/proto"
)

type resrouceInteractor struct {
	resource v1.K8SCustomControllerClient
}

type ResrouceInteractor interface {
	GetResourceSSE(ctx context.Context, w http.ResponseWriter, request models.UsageRequest) error
}

func New(resource v1.K8SCustomControllerClient) ResrouceInteractor {
	return &resrouceInteractor{
		resource: resource,
	}
}

func (r *resrouceInteractor) GetResourceSSE(ctx context.Context, w http.ResponseWriter, request models.UsageRequest) error {
	flusher, _ := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	fmt.Println("ok1")
	go func() {
		for {
			var resp models.UsageResponse
			var usages []models.Usage
			result, err := r.resource.GetResource(ctx, &v1.GetResourceRequest{
				Organization: request.Organization,
				Repository:   request.Repository,
			})
			if err != nil {
				continue
			}

			// Prometheusとかない場合
			if result.IsDisable {
				resp.IsDisable = true
			}
			resources := result.GetResources()
			for _, resource := range resources {
				var usage models.Usage
				usage.PodName = resource.GetPodName()
				for _, metric := range resource.Usages {
					usage = models.Usage{
						CPU:     usage.CPU + int64(metric.Cpu),
						MEM:     usage.MEM + int64(metric.Mem),
						Storage: usage.Storage + int64(metric.Storage),
					}
				}
				usages = append(usages, models.Usage{
					PodName: resource.PodName,
					Branch:  resource.Branch,
					PrID:    resource.PrNumber,
					CPU:     usage.CPU / int64(len(resources)),
					MEM:     usage.MEM / int64(len(resources)),
					Storage: usage.Storage / int64(len(resources)),
				})
			}
			resp.Usages = usages

			data, err := json.Marshal(resp)
			if err != nil {
				continue
			}
			fmt.Fprintf(w, "data: %s\n\n", string(data))

			flusher.Flush()
			time.Sleep(time.Duration(request.Span) * time.Second)
		}
	}()

	<-ctx.Done()
	return nil
}
