package k8sclient

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *k8sClient) CreateConfigMap(ctx context.Context, namespace string, name string, data map[string]string, opts metav1.CreateOptions) error {

	configMap := v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
	_, err := k.client.CoreV1().ConfigMaps(namespace).Create(ctx, &configMap, opts)
	return err
}

func (k *k8sClient) UpdateConfigMap(ctx context.Context, namespace string, name string, data map[string]string, opts metav1.UpdateOptions) error {
	configMap := v1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
	_, err := k.client.CoreV1().ConfigMaps(namespace).Update(ctx, &configMap, opts)
	return err
}

func (k *k8sClient) DeleteConfigMap(ctx context.Context, namespace string, name string, opts metav1.DeleteOptions) error {
	return k.client.CoreV1().ConfigMaps(namespace).Delete(ctx, name, opts)
}
