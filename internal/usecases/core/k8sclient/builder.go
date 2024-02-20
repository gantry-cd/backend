package k8sclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/gantrycd/backend/cmd/config"
	"github.com/gantrycd/backend/internal/utils/random"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type BuilderParams struct {
	Namespace     string
	GitLink       string
	DockerBaseDir string
	DockrFilePath string
	ImageName     string
}

func (k *k8sClient) Builder(ctx context.Context, param BuilderParams, opts ...Option) error {
	o := newOption()

	for _, opt := range opts {
		opt(o)
	}

	name, err := random.RandomString(20)
	if err != nil {
		return err
	}

	// https://github.com/gantrycd/test-repository.git

	// Basic auth

	urls := strings.Split(param.GitLink, "//")
	if len(urls) < 2 {
		return fmt.Errorf("invalid git link")
	}

	gitlinks := fmt.Sprintf("%s//%s:%s@%s", urls[0], config.Config.GitHub.Username, config.Config.GitHub.Password, urls[1])

	pod, err := k.client.CoreV1().Pods(param.Namespace).Create(ctx, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "builder-",
			Labels:       o.labelSelector,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  strings.ToLower(name),
					Image: config.Config.Controller.ImageBuilder,
					Env: []v1.EnvVar{
						{
							Name:  "GIT_REPO",
							Value: gitlinks,
						},
						{
							Name:  "IMAGE_NAME",
							Value: fmt.Sprintf("%s/%s", param.ImageName, param.GitLink),
						},
						{
							Name:  "IMAGE_TAG",
							Value: name,
						},
						{
							Name:  "DOCKER_REGISTRY",
							Value: config.Config.Registry.Host,
						},
						{
							Name:  "DOCKER_REGISTRY_USER",
							Value: config.Config.Registry.User,
						},
						{
							Name:  "DOCKER_REGISTRY_PASSWORD",
							Value: config.Config.Registry.Password,
						},
						{
							Name:  "DOCKER_BASE_DIR",
							Value: param.DockerBaseDir,
						},
						{
							Name:  "DOCKER_FILE_PATH",
							Value: param.DockrFilePath,
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	return k.waitForPod(ctx, param.Namespace, pod.Name)
}

func (k *k8sClient) waitForPod(ctx context.Context, namespace, name string) error {
	for {
		pod, err := k.client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if pod.Status.Phase == v1.PodSucceeded || pod.Status.Phase == v1.PodFailed {
			return nil
		}
	}
}
