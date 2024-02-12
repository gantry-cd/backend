package conf

import "gopkg.in/ini.v1"

// Config はアプリ内の設定データを保持する構造体
type config struct {
	Influx struct {
		Host     string `ini:"host"`
		Port     string `ini:"port"`
		User     string `ini:"user"`
		Password string `ini:"password"`
		Database string `ini:"database"`
	} `ini:"influx"`

	ImageRegistry struct {
		Host     string `ini:"host"`
		User     string `ini:"user"`
		Password string `ini:"password"`
	} `ini:"image_registry"`

	path string
	cfg  *ini.File
}

var Config *config
