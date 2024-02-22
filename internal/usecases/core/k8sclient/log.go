package k8sclient

import (
	v1 "k8s.io/api/core/v1"
	restclient "k8s.io/client-go/rest"
)

func (k *k8sClient) GetLogs(namespace string, podName string, option v1.PodLogOptions) *restclient.Request {
	return k.client.CoreV1().Pods(namespace).GetLogs(podName, &option)
}
