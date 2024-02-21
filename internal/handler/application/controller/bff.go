package controller

import (
	"net/http"

	"github.com/gantrycd/backend/internal/models"
	"github.com/gantrycd/backend/internal/usecases/application/bff"
)

type BffController struct {
	interactor bff.BffInteractor
}

func NewBff(interactor bff.BffInteractor) *BffController {
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
	queries := r.URL.Query()

	if err := bc.interactor.GetRepositoryApps(r.Context(), w, models.GetRepositoryAppsRequest{
		Organization: queries.Get(models.QueryOrganization),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (bc *BffController) GetBranchInfo(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if err := bc.interactor.GetBranchInfo(r.Context(), w, models.GetBranchInfoRequest{
		Organization: queries.Get(models.QueryOrganization),
		Repository:   queries.Get(models.QueryRepository),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
