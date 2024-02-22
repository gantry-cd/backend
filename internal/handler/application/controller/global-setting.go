package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gantrycd/backend/cmd/config"
	"github.com/gantrycd/backend/internal/models"
)

type GlobalSettingController struct {
}

func NewGlobalSettingController() *GlobalSettingController {
	return &GlobalSettingController{}
}

func (gc *GlobalSettingController) GetGlobalGeneralSetting(w http.ResponseWriter, r *http.Request) {
	// TODO: Generalができたら実装する
	if json.NewEncoder(w).Encode(models.GlobalConfigGeneral{}) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (gc *GlobalSettingController) UpdateGlobalGeneralSetting(w http.ResponseWriter, r *http.Request) {
	// TODO: Generalができたら実装する
	if json.NewEncoder(w).Encode(models.GlobalConfigGeneral{}) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (gc *GlobalSettingController) GetGlobalRegistrySetting(w http.ResponseWriter, r *http.Request) {
	if json.NewEncoder(w).Encode(models.GlobalConfigRegistry{
		RegistryHost:     config.Config.Registry.Host,
		RegistryUser:     config.Config.Registry.User,
		RegistryPassword: config.Config.Registry.Password,
	}) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (gc *GlobalSettingController) UpdateGlobalRegistrySetting(w http.ResponseWriter, r *http.Request) {
	// 戦術的として、controllerに問い合わせ、configmapの更新を行う
	// その後、config.LoadEnv()をcontroller,bffともに呼び出し、config.Configを更新する
	// その後、config.Configをjsonで返却する

	if json.NewEncoder(w).Encode(models.GlobalConfigRegistry{
		RegistryHost:     config.Config.Registry.Host,
		RegistryUser:     config.Config.Registry.User,
		RegistryPassword: config.Config.Registry.Password,
	}) != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
