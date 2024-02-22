package controller

import (
	"net/http"

	"github.com/gantrycd/backend/internal/models"
	"github.com/gantrycd/backend/internal/usecases/application/resource"
)

type UsageController struct {
	interactor resource.ResourceInteractor
}

func New(interactor resource.ResourceInteractor) *UsageController {
	return &UsageController{
		interactor: interactor,
	}
}

func (uc *UsageController) Usage(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if err := uc.interactor.GetResource(r.Context(), w, models.UsageRequest{
		Organization: queries.Get(models.QueryOrganization),
		Repository:   queries.Get(models.QueryRepository),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
