package controller

import (
	"context"
	"log"

	coreErr "github.com/gantrycd/backend/internal/error"
	"github.com/gantrycd/backend/internal/usecases/core/k8sclient"
	"github.com/gantrycd/backend/internal/utils"
	v1 "github.com/gantrycd/backend/proto/k8s-controller"
	"google.golang.org/protobuf/types/known/emptypb"
	"k8s.io/client-go/kubernetes"
)

const (
	// NamespaceLabel is the label used to identify the namespace
	CreatedLabel = "created"
	Identity     = "gantrycd"
)

type controller struct {
	v1.UnimplementedK8SCustomControllerServer
	client     *kubernetes.Clientset
	interactor k8sclient.K8SClient
}

func NewController(client *kubernetes.Clientset) v1.K8SCustomControllerServer {
	return &controller{
		client:     client,
		interactor: k8sclient.New(client),
	}
}

func (c *controller) CreateNamespace(ctx context.Context, in *v1.CreateNamespaceRequest) (*v1.CreateNamespaceReply, error) {
	ns, err := c.interactor.CreateNamespace(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &v1.CreateNamespaceReply{
		Name: ns.Name,
	}, nil
}

func (c *controller) ListNamespaces(context.Context, *emptypb.Empty) (*v1.ListNamespacesReply, error) {
	ns, err := c.interactor.ListNamespaces(context.Background(), k8sclient.WithCreatedByLabel(Identity))
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
	return &emptypb.Empty{}, c.interactor.DeleteNamespace(ctx, in.Name)
}

func (c *controller) ApplyDeployment(ctx context.Context, in *v1.CreateDeploymentRequest) (*v1.CreateDeploymentReply, error) {
	dep, err := c.interactor.GetDeployment(ctx, in.Namespace, in.Repository, in.PrNumber)
	if err != nil && err != coreErr.ErrDeploymentsNotFound {
		return nil, err
	}
	// リソースが既に存在している場合は、更新する
	if dep != nil {
		// TODO: Apply Deployment
		return &v1.CreateDeploymentReply{
			Name:      dep.Name,
			Namespace: dep.Namespace,
			Version:   dep.ResourceVersion,
		}, nil
	}

	// リソースが存在しない場合は、新規作成する
	dep, err = c.interactor.CreateDeployment(ctx, in.Namespace, in.PodName, in.Image,
		k8sclient.WithRepositoryLabel(in.Repository),
		k8sclient.WithPrIDLabel(in.PrNumber),
		k8sclient.WithEnvirionmentLabel(k8sclient.EnvPreview),
	)
	if err != nil {
		return nil, err
	}

	return &v1.CreateDeploymentReply{
		Name:      dep.Name,
		Namespace: dep.Namespace,
		Version:   dep.ResourceVersion,
	}, nil
}

func (c *controller) DeleteDeployment(ctx context.Context, in *v1.DeleteDeploymentRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.interactor.DeleteDeployment(ctx, in.Namespace, in.Repository, in.PrNumber)
}

func (c *controller) ListOrgsAnsRepos(ctx context.Context, in *emptypb.Empty) (*v1.ListOrgsAnsReposReply, error) {
	ns, err := c.interactor.ListNamespaces(ctx, k8sclient.WithCreatedByLabel(Identity))
	if err != nil {
		return nil, err
	}
	// n + 1?

	var orgsAndRepos []*v1.OrgsAnsRepos
	for _, n := range ns.Items {

		deps, err := c.interactor.ListDeployments(ctx, n.Name, k8sclient.WithCreatedByLabel(Identity))
		if err != nil {
			return nil, err
		}

		var repos []string
		for _, dep := range deps.Items {
			repo, repoOk := dep.Labels[k8sclient.RepositryLabel]
			// もし、RepoではなくAppの場合は、AppをRepoとして扱う
			// これは、GantryのUI上から, Postgresのようなアプリケーションをデプロイした場合もスコープに入れるため
			if app, appOk := dep.Labels[k8sclient.AppLabel]; appOk {
				repo = app
			}

			if !repoOk {
				continue
			}

			repos = append(repos, repo)
		}

		orgsAndRepos = append(orgsAndRepos, &v1.OrgsAnsRepos{
			Orgs:  n.Name,
			Repos: utils.RemoveDuplocate(repos),
		})

	}
	return &v1.ListOrgsAnsReposReply{
		OrgsAnsRepos: orgsAndRepos,
	}, nil
}
