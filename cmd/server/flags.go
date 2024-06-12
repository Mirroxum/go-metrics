package main

import (
	"flag"
)

var FlagServerAddress string

func parseFlags() {
	flag.StringVar(&FlagServerAddress, "a", "localhost:8080", "HTTP-server address")
	flag.Parse()
	if envRunAddr := AppConfig.ServerAddress; envRunAddr != "" {
		FlagServerAddress = envRunAddr
	}
}
