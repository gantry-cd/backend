package resource

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type k8sResource struct {
	metrics *metrics.Clientset
}

type Resource interface {
	GetLoads(ctx context.Context, namespace string) error
}

func New(metrics *metrics.Clientset) Resource {
	return &k8sResource{
		metrics: metrics,
	}
}

func (r *k8sResource) GetLoads(ctx context.Context, namespace string) error {
	metrics, err := r.metrics.MetricsV1beta1().PodMetricses(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, metric := range metrics.Items {
		for _, container := range metric.Containers {
			cpu := container.Usage.Cpu()
			mem := container.Usage.Memory()
			storage := container.Usage.Storage()

			fmt.Printf("[ namespace: %s , podname: %s ] usage cpu: %v, mem %v, storage %v\n", namespace, container.Name, cpu, mem, storage)
		}
	}
	return nil
}
