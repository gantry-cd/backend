package controller

import (
	"net/http"

	"github.com/gantrycd/backend/internal/usecases/bff"
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
