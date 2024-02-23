package resource

import (
	"context"

	"github.com/gantrycd/backend/internal/usecases/core/k8sclient"
	"github.com/gantrycd/backend/internal/utils/branch"
	v1 "github.com/gantrycd/backend/proto"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type k8sResource struct {
	metrics *metrics.Clientset
}

type Resource interface {
	GetLoads(ctx context.Context, namespace, podName string) (*v1.Resource, error)
}

func New(metrics *metrics.Clientset) Resource {
	return &k8sResource{
		metrics: metrics,
	}
}

func (r *k8sResource) GetLoads(ctx context.Context, namespace, podName string) (*v1.Resource, error) {
	metrics, err := r.metrics.MetricsV1beta1().PodMetricses(namespace).Get(ctx, podName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	branchName, ok := metrics.Labels[k8sclient.BaseBranchLabel]
	if !ok {
		branchName = ""
	}

	branchName, _ = branch.TranspileBranchName(branchName)

	prNumber, ok := metrics.Labels[k8sclient.PullRequestID]
	if !ok {
		prNumber = ""
	}

	resources := &v1.Resource{
		AppName:       metrics.Labels[k8sclient.AppLabel],
		PodName:       podName,
		Branch:        branchName,
		PullRequestId: prNumber,
		Usages:        []*v1.Usage{},
	}

	for _, container := range metrics.Containers {
		cpu := container.Usage.Cpu().MilliValue()
		mem, _ := container.Usage.Memory().AsInt64()
		storage, _ := container.Usage.Storage().AsInt64()

		resources.Usages = append(resources.Usages, &v1.Usage{
			ContainerName: container.Name,
			Cpu:           cpu,
			Mem:           mem,
			Storage:       storage,
		})
	}

	return resources, nil
}
