package k8sclient

import (
	"context"

	v1 "github.com/aura-cd/backend/proto"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *k8sClient) GetPodUsage(ctx context.Context, namespace, podName string) (*v1.Pod, error) {
	metrics, err := k.metrics.MetricsV1beta1().PodMetricses(namespace).Get(ctx, podName, metaV1.GetOptions{})
	if err != nil {
		return nil, err
	}

	pods, err := k.GetPods(ctx, namespace, podName)
	if err != nil {
		return nil, err
	}
	// podNameは一意なので、配列の0番目を取得する
	pod := pods[0]
	var containers []*v1.Container

	for i, container := range metrics.Containers {
		// ここでリソースのリクエスト最小値を取得する
		cpuReq, _ := container.Usage.Cpu().AsInt64()
		memReq, _ := container.Usage.Memory().AsInt64()
		storageReq, _ := container.Usage.Storage().AsInt64()

		// ここでリソースのリクエスト最大値を取得する
		cpuLimit, _ := pod.Spec.Containers[i].Resources.Limits.Cpu().AsInt64()
		memLimit, _ := pod.Spec.Containers[i].Resources.Limits.Memory().AsInt64()
		storageLimit, _ := pod.Spec.Containers[i].Resources.Limits.Storage().AsInt64()

		// ここでリソースのリクエストの使用率を取得する
		cpuUsage, _ := container.Usage.Cpu().AsInt64()
		memUsage, _ := container.Usage.Memory().AsInt64()
		storageUsage, _ := container.Usage.Storage().AsInt64()

		containers = append(containers, &v1.Container{
			Name: container.Name,
			Resource: &v1.Resource{
				Cpu: &v1.CPU{
					Usage:   cpuUsage,
					Request: cpuReq,
					Limit:   cpuLimit,
				},
				Memory: &v1.Memory{
					Usage:   memUsage,
					Request: memReq,
					Limit:   memLimit,
				},
				Storage: &v1.Storage{
					Usage:   storageUsage,
					Request: storageReq,
					Limit:   storageLimit,
				},
			},
		})
	}

	return &v1.Pod{
		Name:       podName,
		Status:     string(pod.Status.Phase),
		Containers: containers,
	}, nil
}
