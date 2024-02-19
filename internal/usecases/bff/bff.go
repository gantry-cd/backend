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
