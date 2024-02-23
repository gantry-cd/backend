package usage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gantrycd/backend/internal/models"
	v1 "github.com/gantrycd/backend/proto"
)

type resrouceInteractor struct {
	resource v1.K8SCustomControllerClient
}

type ResrouceInteractor interface {
	GetResource(ctx context.Context, w http.ResponseWriter, request models.UsageRequest) error
}

func New(resource v1.K8SCustomControllerClient) ResrouceInteractor {
	return &resrouceInteractor{
		resource: resource,
	}
}

func (r *resrouceInteractor) GetResource(ctx context.Context, w http.ResponseWriter, request models.UsageRequest) error {
	resp := models.UsageResponse{
		Organization:   request.Organization,
		DeploymentName: request.DeploymentName,
	}

	result, err := r.resource.GetUsage(ctx, &v1.GetUsageRequest{
		Organization:   request.Organization,
		DeploymentName: request.DeploymentName,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Println(result)
	// Prometheusとかない場合
	if result.GetIsDisable() {
		resp.IsDisable = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			return fmt.Errorf("failed to encode response: %w", err)
		}
		return nil
	}

	resource := result.GetResources()
	var podUsage = make([]models.Usage, len(resource))
	for i, metric := range resource {
		for _, containerUsage := range metric.Usages {
			podUsage[i] = models.Usage{
				CPU:     podUsage[i].CPU + int64(containerUsage.Cpu),
				MEM:     podUsage[i].MEM + int64(containerUsage.Mem),
				Storage: podUsage[i].Storage + int64(containerUsage.Storage),
			}
		}

		podUsage[i] = models.Usage{
			PodName: metric.PodName,
			CPU:     podUsage[i].CPU / int64(len(metric.Usages)),
			MEM:     podUsage[i].MEM / int64(len(metric.Usages)),
			Storage: podUsage[i].Storage / int64(len(metric.Usages)),
		}
	}

	resp.Usages = podUsage

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}
