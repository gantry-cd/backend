package k8sclient

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (k *k8sClient) CreateNamespace(ctx context.Context, name string, opts ...Option) (*corev1.Namespace, error) {
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	k.l.Info("creating namespace", "name", name, CreatedByLabel, o.labelSelector[CreatedByLabel])

	return k.client.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: o.labelSelector,
		},
	}, metav1.CreateOptions{})
}

func (k *k8sClient) ListNamespaces(ctx context.Context, opts ...Option) (*corev1.NamespaceList, error) {
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	return k.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{LabelSelector: labels.Set(o.labelSelector).String()})
}

func (k *k8sClient) DeleteNamespace(ctx context.Context, name string) error {
	return k.client.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
}
