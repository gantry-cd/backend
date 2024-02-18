package controller

import (
	"context"
	"errors"
	"fmt"
	"log"

	coreErr "github.com/gantrycd/backend/internal/error"
	"github.com/gantrycd/backend/internal/usecases/core/k8sclient"
	"github.com/gantrycd/backend/internal/usecases/core/resource"
	"github.com/gantrycd/backend/internal/utils/branch"
	v1 "github.com/gantrycd/backend/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	client  *kubernetes.Clientset
	control k8sclient.K8SClient
	metric  resource.Resource
}

func NewController(
	client *kubernetes.Clientset,
	metric resource.Resource,
) v1.K8SCustomControllerServer {
	return &controller{
		client:  client,
		control: k8sclient.New(client),
		metric:  metric,
	}
}

func (c *controller) CreateNamespace(ctx context.Context, in *v1.CreateNamespaceRequest) (*v1.CreateNamespaceReply, error) {
	ns, err := c.control.CreateNamespace(ctx, in.Name)
	if err != nil {
		return nil, err
	}

	return &v1.CreateNamespaceReply{
		Name: ns.Name,
	}, nil
}

func (c *controller) ListNamespaces(context.Context, *emptypb.Empty) (*v1.ListNamespacesReply, error) {
	ns, err := c.control.ListNamespaces(context.Background(), k8sclient.WithCreatedByLabel(Identity))
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
	return &emptypb.Empty{}, c.control.DeleteNamespace(ctx, in.Name)
}

func (c *controller) ApplyDeployment(ctx context.Context, in *v1.CreateDeploymentRequest) (*v1.CreateDeploymentReply, error) {
	dep, err := c.control.GetDeployment(ctx, k8sclient.GetDeploymentParams{
		Namespace:     in.Namespace,
		Repository:    in.Repository,
		PullRequestID: in.PrNumber,
		Branch:        branch.Transpile1123(in.Branch),
	})
	if err != nil && !errors.Is(err, coreErr.ErrDeploymentsNotFound) {
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
	dep, err = c.control.CreateDeployment(ctx, k8sclient.CreateDeploymentParams{
		Namespace: in.Namespace,
		AppName:   in.AppName,
		Image:     in.Image,
	},
		k8sclient.WithRepositoryLabel(in.Repository),
		k8sclient.WithPrIDLabel(in.PrNumber),
		k8sclient.WithEnvirionmentLabel(k8sclient.EnvPreview),
		k8sclient.WithBaseBranchLabel(branch.Transpile1123(in.Branch)),
	)
	if err != nil {
		return nil, err
	}

	service, err := c.control.CreateNodePortService(ctx, k8sclient.CreateServiceNodePortParams{
		Namespace:   in.Namespace,
		ServiceName: in.AppName,
		TargetPort:  80,
	}, k8sclient.WithRepositoryLabel(in.Repository),
		k8sclient.WithPrIDLabel(in.PrNumber),
		k8sclient.WithEnvirionmentLabel(k8sclient.EnvPreview),
		k8sclient.WithBaseBranchLabel(branch.Transpile1123(in.Branch)),
	)
	if err != nil {
		return nil, err
	}

	fmt.Println(service)

	return &v1.CreateDeploymentReply{
		Name:      dep.Name,
		Namespace: dep.Namespace,
		Version:   dep.ResourceVersion,
	}, nil
}

func (c *controller) DeleteDeployment(ctx context.Context, in *v1.DeleteDeploymentRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.control.DeleteDeployment(ctx, in.Namespace,
		k8sclient.WithRepositoryLabel(in.Repository),
		k8sclient.WithPrIDLabel(in.PrNumber),
		// k8sclient.WithBaseBranchLabel(branch.Transpile1123(in.Branch)),
	)
}

func (c *controller) GetOrgRepos(ctx context.Context, in *v1.GetOrgRepoRequest) (*v1.GetOrgReposReply, error) {
	return c.getOrgRepos(ctx, in.Organization)
}

func (c *controller) GetResource(ctx context.Context, in *v1.GetResourceRequest) (*v1.GetResourceReply, error) {
	resource, err := c.metric.GetLoads(ctx, in.GetOrganization(), in.GetRepository())
	if err != nil {
		return &v1.GetResourceReply{
			Resources: resource,
			IsDisable: false,
		}, nil
	}

	return &v1.GetResourceReply{
		Resources: resource,
		IsDisable: true,
	}, nil
}

func (c *controller) GetAlls(context.Context, *emptypb.Empty) (*v1.GetAllsReply, error) {
	namespaces, err := c.control.ListNamespaces(context.Background(), k8sclient.WithCreatedByLabel(Identity))
	if err != nil {
		return nil, err
	}
	var response *v1.GetAllsReply
	for _, ns := range namespaces.Items {
		deployments, err := c.getOrgRepos(context.Background(), ns.Name)
		if err != nil {
			return nil, err
		}

		response.OrgRepos = append(response.OrgRepos, deployments)
	}

	return response, nil
}

func (c *controller) getOrgRepos(ctx context.Context, organization string) (*v1.GetOrgReposReply, error) {
	deployments, err := c.control.ListDeployments(ctx, organization)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list deployments: %v", err)
	}

	var (
		repos []*v1.Repo
		apps  []*v1.App
	)

	for _, d := range deployments.Items {
		prNumber, prOk := d.Labels[k8sclient.PullRequestID]
		branch, brOk := d.Labels[k8sclient.BaseBranchLabel]
		// PR番号とブランチ名が両方ともない場合はAppとして扱う
		if !prOk && !brOk {
			apps = append(apps, &v1.App{
				AppName: d.Name,
				Version: d.ResourceVersion,
				Image:   d.Spec.Template.Spec.Containers[0].Image,
				Age:     d.CreationTimestamp.String(),
			})
			continue
		}

		// PR番号かブランチ名のどちらかが場合はRepoとして扱う
		repos = append(repos, &v1.Repo{
			RepositoryName: d.Labels[k8sclient.RepositoryLabel],
			PrNumber:       prNumber,
			Branch:         branch,
		})
	}

	return &v1.GetOrgReposReply{
		Organization: organization,
		Repos:        repos,
		Apps:         apps,
	}, nil
}
