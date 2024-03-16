package usage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aura-cd/backend/internal/models"
	v1 "github.com/aura-cd/backend/proto"
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

	// Prometheusとかない場合
	if result.GetIsDisable() {
		resp.IsDisable = true
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			return fmt.Errorf("failed to encode response: %w", err)
		}
		return nil
	}

	resource := result.GetUsages()
	var podUsage = make([]models.Usage, len(resource.Pods))
	for i, metric := range resource.Pods {
		cntainers := metric.GetContainers()
		for _, containerUsage := range cntainers {
			podUsage[i] = models.Usage{
				CPU:     podUsage[i].CPU + int64(containerUsage.Resource.Cpu.Usage),
				MEM:     podUsage[i].MEM + int64(containerUsage.Resource.Memory.Request),
				Storage: podUsage[i].Storage + int64(containerUsage.Resource.Storage.Request),
			}
		}

		podUsage[i] = models.Usage{
			PodName: metric.Name,
			CPU:     podUsage[i].CPU / int64(len(cntainers)),
			MEM:     podUsage[i].MEM / int64(len(cntainers)),
			Storage: podUsage[i].Storage / int64(len(cntainers)),
		}
	}

	resp.Usages = podUsage

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}
