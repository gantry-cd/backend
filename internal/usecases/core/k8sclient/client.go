package k8sclient

type k8sClient struct {
}

type K8SClient interface {
}

func New() K8SClient {
	return &k8sClient{}
}
