package api

import (
	"net/http"

	"github.com/aura-cd/backend/internal/models"
	"github.com/aura-cd/backend/internal/usecases/bff"
)

type BffController struct {
	interactor bff.BffInteractor
}

func New(interactor bff.BffInteractor) *BffController {
	return &BffController{
		interactor: interactor,
	}
}

func (bc *BffController) Home(w http.ResponseWriter, r *http.Request) {
	if err := bc.interactor.GetHome(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (bc *BffController) RepositoryApps(w http.ResponseWriter, r *http.Request) {
	organization := r.PathValue(models.ParamOrganization)

	if err := bc.interactor.GetRepositoryApps(r.Context(), w, models.GetRepositoryAppsRequest{
		Organization: organization,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (bc *BffController) GetRepoBranches(w http.ResponseWriter, r *http.Request) {
	organization := r.PathValue(models.ParamOrganization)
	repo := r.PathValue(models.ParamRepository)
	if organization == "" || repo == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := bc.interactor.GetRepoBranches(r.Context(), w, models.GetRepoBranchesRequest{
		Organization: organization,
		Repository:   repo,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
