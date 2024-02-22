package usage

import (
	"net/http"

	"github.com/gantrycd/backend/internal/models"
	"github.com/gantrycd/backend/internal/usecases/application/usage"
)

type UsageController struct {
	interactor usage.ResrouceInteractor
}

func New(interactor usage.ResrouceInteractor) *UsageController {
	return &UsageController{
		interactor: interactor,
	}
}

func (uc *UsageController) Usage(w http.ResponseWriter, r *http.Request) {

	if err := uc.interactor.GetResource(r.Context(), w, models.UsageRequest{
		Organization: r.PathValue(models.ParamOrganization),
		Pod:          r.PathValue(models.ParamPod),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
