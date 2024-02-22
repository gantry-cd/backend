package router

import (
	"github.com/gantrycd/backend/internal/handler/application/controller"
	"github.com/gantrycd/backend/internal/router/middleware"
)

func (r *router) globalSetting() {
	gc := controller.NewGlobalSettingController()

	r.mux.Handle("GET /global-config/general", middleware.BuildChain(r.middleware.Integrate(gc.GetGlobalGeneralSetting)))
	r.mux.Handle("PUT /global-config/general", middleware.BuildChain(r.middleware.Integrate(gc.UpdateGlobalGeneralSetting)))
	r.mux.Handle("GET /global-config/registry", middleware.BuildChain(r.middleware.Integrate(gc.GetGlobalRegistrySetting)))
	r.mux.Handle("PUT /global-config/registry", middleware.BuildChain(r.middleware.Integrate(gc.UpdateGlobalRegistrySetting)))
}
