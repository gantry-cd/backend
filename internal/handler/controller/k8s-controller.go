package controller

import (
	"context"
	"fmt"
	"log"
	"strconv"

	v1 "github.com/gantrycd/backend/proto/k8s-controller"
	"google.golang.org/protobuf/types/known/emptypb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	applyCoreV1 "k8s.io/client-go/applyconfigurations/core/v1"
	applyMetaV1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CustomController struct {
	v1.UnimplementedK8SCustomControllerServer
	client *kubernetes.Clientset
}

func NewCustomController(client *kubernetes.Clientset) v1.K8SCustomControllerServer {
	return &CustomController{
		client: client,
	}
}

func (c *CustomController) CreateNamespace(ctx context.Context, in *v1.CreateNamespaceRequest) (*v1.CreateNamespaceReply, error) {
	ns, err := c.client.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
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

func (c *CustomController) ListNamespaces(context.Context, *emptypb.Empty) (*v1.ListNamespacesReply, error) {
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

func (c *CustomController) DeleteNamespace(ctx context.Context, in *v1.DeleteNamespaceRequest) (*emptypb.Empty, error) {
	log.Println("deleting namespace", in.Name)
	err := c.client.CoreV1().Namespaces().Delete(ctx, in.Name, metav1.DeleteOptions{})
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *CustomController) ApplyDeployment(ctx context.Context, in *v1.CreateDeploymentRequest) (*v1.CreateDeploymentReply, error) {
	reps, err := strconv.Atoi(in.Replicas)
	if err != nil {
		return nil, err
	}

	// ObjectMeta: metav1.ObjectMeta{
	// 	Name: in.PodName[:63],
	// 	Labels: map[string]string{
	// 		"repository": in.Repository,
	// 		"pr-number":  in.PrNumber,
	// 		"created-by": in.CreatedBy,
	// 	},
	// },
	// Spec: appv1.DeploymentSpec{
	// 	Replicas: ptrint(reps),
	// 	Selector: &metav1.LabelSelector{
	// 		MatchLabels: map[string]string{
	// 			"repository": in.Repository,
	// 			"pr-number":  in.PrNumber,
	// 			"created-by": in.CreatedBy,
	// 		},
	// 	},
	// 	Template: corev1.PodTemplateSpec{
	// 		ObjectMeta: metav1.ObjectMeta{
	// 			Name: in.PodName[:63],
	// 			Labels: map[string]string{
	// 				"repository": in.Repository,
	// 				"pr-number":  in.PrNumber,
	// 				"created-by": in.CreatedBy,
	// 			},
	// 		},
	// 		Spec: corev1.PodSpec{
	// 			Containers: []corev1.Container{
	// 				{
	// 					Name:  in.PodName[:63],
	// 					Image: in.Image,
	// 				},
	// 			},
	// 		},
	// 	},
	// },

	dep, err := c.client.AppsV1().Deployments(in.Namespace).Apply(ctx, &appsv1.DeploymentApplyConfiguration{
		Spec: &appsv1.DeploymentSpecApplyConfiguration{
			Replicas: ptrint(reps),
			Selector: &applyMetaV1.LabelSelectorApplyConfiguration{
				MatchLabels: map[string]string{
					"repository": in.Repository,
					"pr-number":  in.PrNumber,
					"created-by": in.CreatedBy,
				},
			},
			Template: &applyCoreV1.PodTemplateSpecApplyConfiguration{
				Spec: &applyCoreV1.PodSpecApplyConfiguration{
					Containers: []applyCoreV1.ContainerApplyConfiguration{
						{
							Name:  &in.PodName,
							Image: &in.Image,
						},
					},
				},
			},
		},
	}, metav1.ApplyOptions{})

	if err != nil {
		return nil, err
	}

	return &v1.CreateDeploymentReply{
		Name: dep.Name,
	}, nil
}

func (c *CustomController) DeleteDeployment(ctx context.Context, in *v1.DeleteDeploymentRequest) (*emptypb.Empty, error) {
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
