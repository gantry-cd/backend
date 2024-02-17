package resource

import (
	"context"

	"github.com/gantrycd/backend/internal/usecases/core/k8sclient"
	v1 "github.com/gantrycd/backend/proto"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type k8sResource struct {
	metrics *metrics.Clientset
}

type Resource interface {
	GetLoads(ctx context.Context, namespace, repository string) ([]*v1.Resource, error)
}

func New(metrics *metrics.Clientset) Resource {
	return &k8sResource{
		metrics: metrics,
	}
}

func (r *k8sResource) GetLoads(ctx context.Context, namespace, repository string) ([]*v1.Resource, error) {
	metrics, err := r.metrics.MetricsV1beta1().PodMetricses(namespace).List(ctx, metaV1.ListOptions{
		LabelSelector: labels.Set(map[string]string{
			k8sclient.RepositryLabel: repository,
		}).String(),
	})
	if err != nil {
		return nil, err
	}
	var resources []*v1.Resource

	for _, metric := range metrics.Items {
		var usages []*v1.Usage
		for _, container := range metric.Containers {
			cpu, _ := container.Usage.Cpu().AsInt64()
			mem, _ := container.Usage.Memory().AsInt64()
			storage, _ := container.Usage.Storage().AsInt64()

			usages = append(usages, &v1.Usage{
				ContainerName: container.Name,
				Cpu:           cpu,
				Mem:           mem,
				Storage:       storage,
			})
		}

		resources = append(resources, &v1.Resource{
			AppName: metric.Labels[k8sclient.AppLabel],
			PodName: metric.Name,
			Usages:  usages,
		})
	}

	return resources, nil
}
