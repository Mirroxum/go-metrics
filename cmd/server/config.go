package main

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string `env:"ADDRESS"`
	LogLevel      string `env:"LOG_LEVEL"`
}

var AppConfig Config

func loadConfig() {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		panic(err)
	}
	AppConfig = config
}
