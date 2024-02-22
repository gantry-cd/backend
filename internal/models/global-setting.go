package models

const (
	QueryRegistryHost     string = "registryHost"
	QueryRegistryUser     string = "registryUser"
	QueryRegistryPassword string = "registryPassword"
)

type GlobalConfigRegistry struct {
	RegistryHost     string `json:"registryHost"`
	RegistryUser     string `json:"registryUser"`
	RegistryPassword string `json:"registryPassword"`
}

type GlobalConfigGeneral struct{}
