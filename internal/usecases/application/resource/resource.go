package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gantrycd/backend/internal/models"
	v1 "github.com/gantrycd/backend/proto"
)

type resourceInteractor struct {
	resource v1.K8SCustomControllerClient
}

type ResourceInteractor interface {
	GetResource(ctx context.Context, w http.ResponseWriter, request models.UsageRequest) error
}

func New(resource v1.K8SCustomControllerClient) ResourceInteractor {
	return &resourceInteractor{
		resource: resource,
	}
}

func (r *resourceInteractor) GetResource(ctx context.Context, w http.ResponseWriter, request models.UsageRequest) error {
	var resp models.UsageResponse
	var usages []models.Usage
	result, err := r.resource.GetResource(ctx, &v1.GetResourceRequest{
		Organization: request.Organization,
		Repository:   request.Repository,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
			PrID:    resource.PullRequestId,
			CPU:     usage.CPU / int64(len(resources)),
			MEM:     usage.MEM / int64(len(resources)),
			Storage: usage.Storage / int64(len(resources)),
		})
	}
	resp.Usages = usages

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}
