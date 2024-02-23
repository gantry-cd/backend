package k8sclient

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *k8sClient) CreateConfigMap(ctx context.Context, configMap v1.ConfigMap, opts metav1.CreateOptions) error {
	_, err := k.client.CoreV1().ConfigMaps(configMap.Namespace).Create(ctx, &configMap, opts)
	return err
}

func (k *k8sClient) UpdateConfigMap(ctx context.Context, configMap v1.ConfigMap, opts metav1.UpdateOptions) error {
	_, err := k.client.CoreV1().ConfigMaps(configMap.Namespace).Update(ctx, &configMap, opts)
	return err
}

func (k *k8sClient) DeleteConfigMap(ctx context.Context, namespace string, name string, opts metav1.DeleteOptions) error {
	return k.client.CoreV1().ConfigMaps(namespace).Delete(ctx, name, opts)
}
