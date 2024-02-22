package log

import (
	"net/http"

	"github.com/gantrycd/backend/internal/models"
	"github.com/gantrycd/backend/internal/usecases/application/logger"
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
		Organization: r.PathValue(models.QueryOrganization),
		Pod:          r.PathValue(models.QueryPod),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
