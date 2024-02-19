package models

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type PullRequestConfig struct {
	BuildFilePath string      `json:"build_file_path" conf:"BUILD_FILE_PATH"`
	ExposePort    []string    `json:"expose_port" conf:"EXPOSE_PORT"`
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

func ParseConfig(configStr string) (*PullRequestConfig, error) {
	config := new(PullRequestConfig)

	// 構造体のフィールドに直接アクセス
	for _, line := range strings.Split(configStr, "\n") {

		// コメント行は無視
		if strings.HasPrefix(line, "#") {
			continue
		}

		// キーと値に分割
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}

		t := reflect.TypeOf(*config)
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			tag := field.Tag.Get("conf")
			if tag == "" {
				continue
			}

			// キーが一致するフィールドを探す
			if tag == parts[0] {
				// 値をセット
				if err := setStructField(config, field, parts[1]); err != nil {
					return nil, err
				}
			}
		}

	}

	return config, nil
}
func setStructField(config interface{}, field reflect.StructField, value string) error {
	fieldValue := reflect.ValueOf(config).Elem().FieldByName(field.Name)
	if !fieldValue.CanSet() {
		return fmt.Errorf("field %s is not settable", field.Name)
	}

	switch fieldValue.Kind() {
	case reflect.String:
		fieldValue.SetString(value)
	case reflect.Int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		fieldValue.SetInt(int64(intValue))
	case reflect.Slice:
		elemType := fieldValue.Type().Elem()
		if elemType.Kind() == reflect.String {
			fieldValue.Set(reflect.MakeSlice(fieldValue.Type(), 0, 1))
			fieldValue.Index(0).SetString(value)
		}
	default:
		return fmt.Errorf("unsupported field type: %s", fieldValue.Type())
	}

	return nil
}
