package k8sclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aura-cd/backend/cmd/config"
	"github.com/aura-cd/backend/internal/utils"
	"github.com/aura-cd/backend/internal/utils/random"
	"github.com/aura-cd/backend/internal/utils/url"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	BatchV1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Image Builder Environment variables
	EnvGitRepo        = "GIT_REPO"
	EnvGitBranch      = "GIT_BRANCH"
	EnvImageName      = "IMAGE_NAME"
	EnvImageTag       = "IMAGE_TAG"
	EnvDockerBaseDir  = "DOCKER_BASE_DIR"
	EnvDockerFilePath = "DOCKER_FILE_PATH"
	EnvDockerRegistry = "DOCKER_REGISTRY"
	EnvDockerUser     = "DOCKER_REGISTRY_USER"
	EnvDockerPassword = "DOCKER_REGISTRY_PASSWORD"
)

type ImageBuilderParams struct {
	Repository     string
	GitRepo        string
	GitBranch      string
	ImageName      string
	ImageTag       string
	DockerBaseDir  string
	DockerFilePath string
	Token          string
}

// imageBuilderEnv はイメージビルダーの環境変数とimageを生成する .
func imageBuilderEnv(param ImageBuilderParams) ([]v1.EnvVar, string) {
	tag, _ := random.RandomString(20)
	image := fmt.Sprintf("%s/%s/%s", config.Config.Registry.Host, config.Config.Application.ApplicationName, param.Repository)

	return []v1.EnvVar{
		toEnvVar(EnvGitRepo, url.IncludeBasicAuth(param.GitRepo, config.Config.GitHub.Username, param.Token)),
		toEnvVar(EnvGitBranch, param.GitBranch),
		toEnvVar(EnvImageName, image),
		toEnvVar(EnvImageTag, tag),
		toEnvVar(EnvDockerBaseDir, param.DockerBaseDir),
		toEnvVar(EnvDockerFilePath, param.DockerFilePath),
		toEnvVar(EnvDockerRegistry, config.Config.Registry.Host),
		toEnvVar(EnvDockerUser, config.Config.Registry.User),
		toEnvVar(EnvDockerPassword, config.Config.Registry.Password),
	}, fmt.Sprintf("%s:%s", image, tag)
}

const (
	// DefaultBackoffLimit はジョブのバックオフリミットのデフォルト値 .
	DefaultBackoffLimit = 3
)

type BuilderParams struct {
	Namespace string
	Branch    string

	BuilderParam ImageBuilderParams
}

func (k *k8sClient) Builder(ctx context.Context, param BuilderParams, opts ...Option) (string, error) {
	o := newOption()

	for _, opt := range opts {
		opt(o)
	}

	name, err := random.RandomString(20)
	if err != nil {
		return "", err
	}

	envVar, image := imageBuilderEnv(param.BuilderParam)

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
							Env:   envVar,
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
		return "", status.Errorf(codes.Internal, err.Error())
	}

	if err := k.waitForJob(ctx, param.Namespace, job.Name); err != nil {
		return "", err
	}
	fmt.Println(image)
	return image, nil
}

func (k *k8sClient) waitForJob(ctx context.Context, namespace, name string) error {
	for {
		job, err := k.client.BatchV1().Jobs(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}

		// jobが成功したらjob消して終了
		if job.Status.Succeeded > 0 {
			return nil
		}

		if job.Status.Failed >= 3 {
			return status.Errorf(codes.Internal, "job failed")
		}

		time.Sleep(2 * time.Second)
	}
}
