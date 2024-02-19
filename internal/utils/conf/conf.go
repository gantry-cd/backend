package conf

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/gantrycd/backend/internal/models"
)

const (
	configMapField = "[config-map]"
)

// LoadConf は指定されたファイルパスの設定ファイルを読み込み、指定された構造体にマッピングする
// path: 設定ファイルのパス
// model: 設定ファイルをマッピングする構造体のポインタ
// if err := conf.LoadConf("path/to/config.ini", &conf.Config); err != nil { ... }
func LoadConf(conf string) (*models.PullRequestConfig, error) {
	var (
		config models.PullRequestConfig
		prefix string
	)

	for _, line := range strings.Split(conf, "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, configMapField) {
			prefix = configMapField
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}

		if err := setValue(&config, prefix, parts[0], parts[1]); err != nil {
			return nil, err
		}
	}
	return &config, nil
}

func setValue(model *models.PullRequestConfig, prefix, key, value string) error {
	configType := reflect.TypeOf(*model)
	configValue := reflect.ValueOf(model)
	for i := 0; i < configType.NumField(); i++ {
		switch prefix {
		case "":
			if configType.Field(i).Tag.Get("conf") == key {
				if configType.Field(i).Type.Kind() == reflect.Slice {
					if configType.Field(i).Type.Elem().Kind() == reflect.Int32 {
						v, err := strconv.Atoi(value)
						if err != nil {
							return err
						}

						configValue.Elem().Field(i).Set(reflect.Append(configValue.Elem().Field(i), reflect.ValueOf(int32(v))))
						continue
					}

					if configType.Field(i).Type.Elem().Kind() == reflect.String {
						configValue.Elem().Field(i).Set(reflect.Append(configValue.Elem().Field(i), reflect.ValueOf(value)))
						continue
					}
				}
				log.Println(value)
				configValue.Elem().Field(i).SetString(value)
			}
		case configMapField:
			model.ConfigMaps = append(model.ConfigMaps, models.SetConfigMap(key, value))
			return nil
		default:
			continue
		}
	}
	return nil
}
