package k8sclient

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *k8sClient) GetReplicaSet(ctx context.Context, namespace string, prefix string) (*appsv1.ReplicaSet, error) {
	rs, err := k.client.AppsV1().ReplicaSets(namespace).Get(ctx, prefix, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return rs, nil
}
