package k8sclient

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *k8sClient) GetPods(ctx context.Context, namespace, prefix string) ([]*corev1.Pod, error) {
	pods, err := k.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []*corev1.Pod
	for _, pod := range pods.Items {
		if prefix != "" && !strings.HasPrefix(pod.Name, prefix) {
			continue
		}

		result = append(result, &pod)
	}

	return result, nil
}
