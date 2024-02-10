package k8s

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type k8sClient struct {
	kubeconfig *string
}

type K8sClient interface {
	ConnectDynamic(masterURL string) (*dynamic.DynamicClient, error)
	ConnectTyped(masterURL string) (*kubernetes.Clientset, error)
	ConnectMetrics(masterURL string) (*metrics.Clientset, error)
}

func New(kubeConfig *string) K8sClient {
	return &k8sClient{
		kubeconfig: kubeConfig,
	}
}

func (k *k8sClient) ConnectDynamic(masterURL string) (*dynamic.DynamicClient, error) {
	config, err := clientcmd.BuildConfigFromFlags(masterURL, *k.kubeconfig)
	if err != nil {
		return nil, err
	}

	return dynamic.NewForConfig(config)
}

func (k *k8sClient) ConnectTyped(masterURL string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags(masterURL, *k.kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func (k *k8sClient) ConnectMetrics(masterURL string) (*metrics.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags(masterURL, *k.kubeconfig)
	if err != nil {
		return nil, err
	}

	return metrics.NewForConfig(config)
}
