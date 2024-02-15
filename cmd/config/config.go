package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

func LoadEnv(path ...string) {
	if path != nil {
		if err := godotenv.Load(path...); err != nil {
			panic(err)
		}
	}

	config := &config{}

	if err := env.Parse(&config.Server); err != nil {
		panic(err)
	}

	if err := env.Parse(&config.KeyCloak); err != nil {
		panic(err)
	}

	fmt.Println(config)
	Config = config
}
