package controller

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"

	v1 "github.com/gantrycd/backend/proto/k8s-controller"
	"google.golang.org/protobuf/types/known/emptypb"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	// NamespaceLabel is the label used to identify the namespace
	CreatedLabel = "created"
	Identity     = "gantrycd"
)

type controller struct {
	v1.UnimplementedK8SCustomControllerServer
	client *kubernetes.Clientset
}

func NewController(client *kubernetes.Clientset) v1.K8SCustomControllerServer {
	return &controller{
		client: client,
	}
}

func (c *controller) CreateNamespace(ctx context.Context, in *v1.CreateNamespaceRequest) (*v1.CreateNamespaceReply, error) {
	ns, err := c.client.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				CreatedLabel: Identity,
			},
			Name: in.Name,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return &v1.CreateNamespaceReply{
		Name: ns.Name,
	}, nil
}

func (c *controller) ListNamespaces(context.Context, *emptypb.Empty) (*v1.ListNamespacesReply, error) {
	ns, err := c.client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, n := range ns.Items {
		log.Println(n.Name)
		names = append(names, n.Name)
	}

	return &v1.ListNamespacesReply{
		Names: names,
	}, nil
}

func (c *controller) DeleteNamespace(ctx context.Context, in *v1.DeleteNamespaceRequest) (*emptypb.Empty, error) {
	log.Println("deleting namespace", in.Name)
	err := c.client.CoreV1().Namespaces().Delete(ctx, in.Name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *controller) ApplyDeployment(ctx context.Context, in *v1.CreateDeploymentRequest) (*v1.CreateDeploymentReply, error) {
	reps, err := strconv.Atoi(in.Replicas)
	if err != nil {
		return nil, err
	}

	if len(in.CreatedBy) == 0 {
		in.CreatedBy = Identity
	}

	in.PodName = fmt.Sprintf("%s-%s", in.Repository, in.PrNumber)

	dep, err := c.client.AppsV1().Deployments(in.Namespace).Create(ctx, &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: in.PodName,
			Labels: map[string]string{
				"repository": in.Repository,
				"pr-number":  in.PrNumber,
				"created-by": in.CreatedBy,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: ptrint(reps),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"repository": in.Repository,
					"pr-number":  in.PrNumber,
					"created-by": in.CreatedBy,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: in.PodName,
					Labels: map[string]string{
						"repository": in.Repository,
						"pr-number":  in.PrNumber,
						"created-by": in.CreatedBy,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  in.PodName,
							Image: in.Image,
						},
					},
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return &v1.CreateDeploymentReply{
		Name: dep.Name,
	}, nil
}

func (c *controller) DeleteDeployment(ctx context.Context, in *v1.DeleteDeploymentRequest) (*emptypb.Empty, error) {
	deps, err := c.client.AppsV1().Deployments(in.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: fmt.Sprintf("repository=%s,pr-number=%s", in.Repository, in.PrNumber),
	})

	if err != nil {
		return nil, err
	}

	log.Println("deleting deployments", deps.Items)

	for _, dep := range deps.Items {
		log.Println("deleting deployment", dep.Name)
		// delete deployment
		err := c.client.AppsV1().Deployments(in.Namespace).Delete(ctx, dep.Name, metav1.DeleteOptions{})
		if err != nil {
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}

func ptrint(i int) *int32 {
	j := int32(i)
	return &j
}

func (c *controller) ListOrgsAnsRepos(context.Context, *emptypb.Empty) (*v1.ListOrgsAnsReposReply, error) {
	ns, err := c.client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", CreatedLabel, Identity),
	})

	if err != nil {
		return nil, err
	}
	// n + 1?

	var orgsAndRepos []*v1.OrgsAnsRepos
	for _, n := range ns.Items {

		deps, err := c.client.AppsV1().Deployments(n.Name).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		var repos []string
		for _, dep := range deps.Items {
			repo, repoOk := dep.Labels["repository"]
			// もし、RepoではなくAppの場合は、AppをRepoとして扱う
			// これは、GantryのUI上から, Postgresのようなアプリケーションをデプロイした場合もスコープに入れるため
			if app, appOk := dep.Labels["app"]; appOk {
				repo = app
			}

			if !repoOk {
				continue
			}

			repos = append(repos, repo)
		}

		orgsAndRepos = append(orgsAndRepos, &v1.OrgsAnsRepos{
			Orgs:  n.Name,
			Repos: removeDuplocate(repos),
		})

	}
	return &v1.ListOrgsAnsReposReply{
		OrgsAnsRepos: orgsAndRepos,
	}, nil
}

// removeDuplocate removes duplicate elements from a slice
func removeDuplocate[T any](s []T) []T {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if reflect.DeepEqual(s[i], s[j]) {
				s = append(s[:j], s[j+1:]...)
				j--
			}
		}
	}
	return s
}
