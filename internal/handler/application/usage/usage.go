package usage

import (
	"net/http"

	"github.com/aura-cd/backend/internal/models"
	"github.com/aura-cd/backend/internal/usecases/application/usage"
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
	organization := r.PathValue(models.ParamOrganization)
	deploymentName := r.PathValue(models.ParamDeploymentName)
	if organization == "" || deploymentName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := uc.interactor.GetResource(r.Context(), w, models.UsageRequest{
		Organization:   organization,
		DeploymentName: deploymentName,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
