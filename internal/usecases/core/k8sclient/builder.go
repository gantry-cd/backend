package k8sclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/gantrycd/backend/cmd/config"
	"github.com/gantrycd/backend/internal/utils"
	"github.com/gantrycd/backend/internal/utils/random"
	BatchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	DefaultBackoffLimit = 3
)

type BuilderParams struct {
	Namespace     string
	Repository    string
	Branch        string
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

	job, err := k.client.BatchV1().Jobs(param.Namespace).Create(ctx, &BatchV1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "build-job-",
		},
		Spec: BatchV1.JobSpec{
			Template: v1.PodTemplateSpec{
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
									Value: fmt.Sprintf("%s/%s", param.ImageName, param.Repository),
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
									Name:  "GIT_BRANCH",
									Value: param.Branch,
								},
								{
									Name:  "DOCKER_FILE_PATH",
									Value: param.DockrFilePath,
								},
							},
							SecurityContext: &v1.SecurityContext{
								Privileged: utils.ToPtr(true),
							},
						},
					},
					RestartPolicy: "Never",
				},
			},
			BackoffLimit: utils.ToPtr(int32(DefaultBackoffLimit)),
			Completions:  utils.ToPtr(int32(1)),
			Parallelism:  utils.ToPtr(int32(1)),
		},
	}, metav1.CreateOptions{})

	if err != nil {
		return err
	}

	return k.waitForJob(ctx, param.Namespace, job.Name)
}

func (k *k8sClient) waitForJob(ctx context.Context, namespace, name string) error {
	for {
		job, err := k.client.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if job.Status.Succeeded > 0 {
			return nil
		}

		if job.Status.Failed == DefaultBackoffLimit {
			return fmt.Errorf("job failed")
		}
	}
}
