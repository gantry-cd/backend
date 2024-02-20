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
		configField := configType.Field(i)
		configValue := configValue.Elem().Field(i)
		switch prefix {
		case "":
			if configType.Field(i).Tag.Get("conf") == key {
				log.Println(value)
				if err := set(configField, configValue, value); err != nil {
					return err
				}
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

func set(confType reflect.StructField, confValue reflect.Value, value string) error {

	switch confType.Type.Kind() {
	case reflect.String:
		confValue.SetString(value)

	case reflect.Int:
		v, err := strconv.Atoi(strings.ReplaceAll(value, " ", ""))
		if err != nil {
			return err
		}
		confValue.SetInt(int64(v))

	case reflect.Int32:
		v, err := strconv.Atoi(strings.ReplaceAll(value, " ", ""))
		if err != nil {
			return err
		}
		confValue.SetInt(int64(v))

	case reflect.Int64:
		v, err := strconv.Atoi(strings.ReplaceAll(value, " ", ""))
		if err != nil {
			return err
		}
		confValue.SetInt(int64(v))

	case reflect.Slice:
		switch confType.Type.Elem().Kind() {
		case reflect.String:
			confValue.Set(reflect.Append(confValue, reflect.ValueOf(value)))

		case reflect.Int:
			v, err := strconv.Atoi(strings.ReplaceAll(value, " ", ""))
			if err != nil {
				return err
			}
			confValue.Set(reflect.Append(confValue, reflect.ValueOf(v)))

		case reflect.Int32:
			v, err := strconv.Atoi(strings.ReplaceAll(value, " ", ""))
			if err != nil {
				return err
			}
			confValue.Set(reflect.Append(confValue, reflect.ValueOf(int32(v))))

		case reflect.Int64:
			v, err := strconv.Atoi(strings.ReplaceAll(value, " ", ""))
			if err != nil {
				return err
			}
			confValue.Set(reflect.Append(confValue, reflect.ValueOf(int64(v))))
		}

	}

	return nil
}
