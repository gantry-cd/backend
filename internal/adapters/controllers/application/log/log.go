package log

import (
	"net/http"

	"github.com/aura-cd/backend/internal/models"
	"github.com/aura-cd/backend/internal/usecases/application/logger"
)

type Controller struct {
	interactor logger.LogInteractor
}

func New(interactor logger.LogInteractor) *Controller {
	return &Controller{
		interactor: interactor,
	}
}

func (uc *Controller) Log(w http.ResponseWriter, r *http.Request) {
	if err := uc.interactor.GetLogStream(r.Context(), w, models.PodLogRequest{
		Organization: r.PathValue(models.ParamOrganization),
		Pod:          r.PathValue(models.ParamPod),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
