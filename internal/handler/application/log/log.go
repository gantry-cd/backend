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
	queries := r.URL.Query()

	if err := uc.interactor.GetLogStream(r.Context(), w, models.PodLogRequest{
		Organization: queries.Get(models.QueryOrganization),
		Repository:   queries.Get(models.QueryRepository),
		Pull:         queries.Get(models.QueryPull),
		Pod:          queries.Get(models.QueryPod),
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
