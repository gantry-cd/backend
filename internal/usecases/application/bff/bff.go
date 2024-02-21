package bff

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gantrycd/backend/internal/models"
	v1 "github.com/gantrycd/backend/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type bffInteractor struct {
	resource v1.K8SCustomControllerClient
}

type BffInteractor interface {
	GetHome(ctx context.Context, w http.ResponseWriter) error
	GetRepositoryApps(ctx context.Context, w http.ResponseWriter, request models.GetRepositoryAppsRequest) error
}

func NewBff(resource v1.K8SCustomControllerClient) BffInteractor {
	return &bffInteractor{
		resource: resource,
	}
}

func (b *bffInteractor) GetHome(ctx context.Context, w http.ResponseWriter) error {
	var orgInfos []models.OrganizationInfos
	result, err := b.resource.GetAlls(ctx, &emptypb.Empty{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	results := result.GetOrganizationInfos()
	if len(results) == 0 {
		http.Error(w, "No data", http.StatusNotFound)
		return nil
	}
	for _, result := range results {
		repos := result.GetRepositories()
		var resultRepo []string
		for _, repo := range repos {
			resultRepo = append(resultRepo, repo.Name)
		}
		orgInfos = append(orgInfos, models.OrganizationInfos{
			Organization: result.GetOrganization(),
			Repositories: resultRepo,
		})
	}
	if err := json.NewEncoder(w).Encode(models.HomeResponse{OrganizationInfos: orgInfos}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (b *bffInteractor) GetRepositoryApps(ctx context.Context, w http.ResponseWriter, request models.GetRepositoryAppsRequest) error {
	var resApps []models.Apps
	var resRepos []models.Repositories
	var repoIndex []string
	repoCount := make(map[string]int)
	result, err := b.resource.GetOrgRepos(ctx, &v1.GetOrgRepoRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	repos := result.GetRepositories()
	apps := result.GetApplications()

	for _, repo := range repos {
		if !containsValue(repoIndex, repo.Name) {
			repoIndex = append(repoIndex, repo.Name)
			repoCount[repo.Name] = 1
		} else {
			repoCount[repo.Name] += 1
		}
	}

	for _, repoName := range repoIndex {
		resRepos = append(resRepos, models.Repositories{
			Repository: repoName,
			Deployment: int32(repoCount[repoName]),
		})
	}

	for _, app := range apps {
		resApps = append(resApps, models.Apps{
			Name:    app.Name,
			Status:  app.Status,
			Version: app.Version,
			Age:     app.Age,
		})
	}
	if err := json.NewEncoder(w).Encode(models.GetRepositoryAppsResponse{Repositories: resRepos, Apps: resApps}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func containsValue(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func (b *bffInteractor) GetBranchInfo(ctx context.Context, w http.ResponseWriter, request models.GetBranchInfoRequest) error {
	var resBranches models.GerBranchInfoResponse
	result, err := b.resource.GetBranchInfo(ctx, &v1.GetBranchInfoRequest{
		Organization: request.Organization,
		Repository:   request.Repository,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	branches := result.GetBranches()
	for _, branch := range branches {
		resBranches = append(resBranches, models.Branches{})
	}
	if err := json.NewEncoder(w).Encode(models.GetBranchInfoResponse{Branches: resBranches}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return nil
}
