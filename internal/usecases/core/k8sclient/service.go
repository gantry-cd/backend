package k8sclient

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type CreateServiceNodePortParams struct {
	Namespace   string
	ServiceName string
	TargetPort  []int32
}

func (k *k8sClient) CreateNodePortService(ctx context.Context, param CreateServiceNodePortParams, opts ...Option) (*corev1.Service, error) {
	o := newOption()
	for _, opt := range opts {
		opt(o)
	}

	var expose []corev1.ServicePort
	for _, port := range param.TargetPort {
		expose = append(expose, corev1.ServicePort{
			Port: port,
		})

	}

	return k.client.CoreV1().Services(param.Namespace).Create(ctx, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", param.ServiceName),
			Labels:       o.labelSelector,
		},
		Spec: corev1.ServiceSpec{
			Ports:    expose,
			Selector: o.labelSelector,
			Type:     corev1.ServiceTypeNodePort,
		},
	}, metav1.CreateOptions{})
}

func (k *k8sClient) DeleteService(ctx context.Context, namespace string, opts ...Option) error {
	o := newOption()

	for _, opt := range opts {
		opt(o)
	}

	services, err := k.client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.Set(o.labelSelector).String(),
	})

	if err != nil {
		return err
	}

	for _, service := range services.Items {
		if err := k.client.CoreV1().Services(namespace).Delete(ctx, service.Name, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}

	return nil
}
