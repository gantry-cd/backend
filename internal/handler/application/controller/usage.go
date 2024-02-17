package controller

import (
	"net/http"
	"strconv"

	"github.com/gantrycd/backend/internal/models"
	"github.com/gantrycd/backend/internal/usecases/application/resource"
)

type UsageController struct {
	interactor resource.ResrouceInteractor
}

func New(interactor resource.ResrouceInteractor) *UsageController {
	return &UsageController{
		interactor: interactor,
	}
}

func (uc *UsageController) Usage(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	span, err := strconv.Atoi(queries.Get(models.QuerySpan))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := uc.interactor.GetResourceSSE(r.Context(), w, models.UsageRequest{
		Organization: queries.Get(models.QueryOrganization),
		Repository:   queries.Get(models.QueryRepository),
		Span:         span,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
