package models

type PullRequestConfig struct {
	BuildFilePath string      `json:"build_file_path" conf:"BUILD_FILE_PATH"`
	Replicas      int32       `json:"replicas" conf:"REPLICAS"`
	ExposePort    []int32     `json:"expose_port" conf:"EXPOSE_PORT"`
	ConfigMaps    []ConfigMap `json:"config_map" conf:"[config-map]"`
}

type ConfigMap struct {
	Name  string `json:"name" conf:"key"`
	Value string `json:"value" conf:"value"`
}

func SetConfigMap(key, value string) ConfigMap {
	return ConfigMap{
		Name:  key,
		Value: value,
	}
}
