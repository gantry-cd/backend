package usage

import (
	"context"
	"encoding/json"
	"fmt"
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
	var usages []models.Usage
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

	resource := result.GetResources()
	var usage models.Usage
	usage.PodName = resource.GetPodName()
	for _, metric := range resource.GetUsages() {
		usage = models.Usage{
			CPU:     usage.CPU + int64(metric.Cpu),
			MEM:     usage.MEM + int64(metric.Mem),
			Storage: usage.Storage + int64(metric.Storage),
		}
	}
	usages = append(usages, models.Usage{
		PodName: resource.PodName,
		CPU:     usage.CPU / int64(len(resource.Usages)),
		MEM:     usage.MEM / int64(len(resource.Usages)),
		Storage: usage.Storage / int64(len(resource.Usages)),
	})
	resp.Usages = usages

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}
