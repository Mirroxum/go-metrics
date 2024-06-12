package main

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress  string `env:"ADDRESS"`
	ReportInterval int64  `env:"REPORT_INTERVAL"`
	PollInterval   int64  `env:"POLL_INTERVAL"`
}

var AppConfig Config

func loadConfig() {
	config := Config{}
	if err := env.Parse(&config); err != nil {
		panic(err)
	}
	AppConfig = config
}
