package config

import (
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

	if err := env.Parse(&config.Bff); err != nil {
		panic(err)
	}
	if err := env.Parse(&config.Controller); err != nil {
		panic(err)
	}
	if err := env.Parse(&config.GitHub); err != nil {
		panic(err)
	}

	if err := env.Parse(&config.KeyCloak); err != nil {
		panic(err)
	}

	if err := env.Parse(&config.Registry); err != nil {
		panic(err)
	}

	Config = config
}
