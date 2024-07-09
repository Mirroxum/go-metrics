package main

import (
	"flag"
)

var (
	FlagServerAddress string
	FlagLogLevel      string
)

func parseFlags() {
	flag.StringVar(&FlagServerAddress, "a", "localhost:8080", "HTTP-server address")
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")
	flag.Parse()
	if envRunAddr := AppConfig.ServerAddress; envRunAddr != "" {
		FlagServerAddress = envRunAddr
	}
	if envLogLevel := AppConfig.LogLevel; envLogLevel != "" {
		FlagLogLevel = envLogLevel
	}
}
