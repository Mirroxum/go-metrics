package main

import (
	"flag"
)

var (
	FlagServerAddress  string
	FlagReportInterval int64
	FlagPollInterval   int64
)

func parseFlags() {
	flag.StringVar(&FlagServerAddress, "a", "localhost:8080", "HTTP-server address")
	flag.Int64Var(&FlagReportInterval, "r", 10, "Frequency of sending metrics to the server")
	flag.Int64Var(&FlagPollInterval, "p", 2, "Frequency of polling metrics from the runtime package")
	flag.Parse()
	if envRunAddr := AppConfig.ServerAddress; envRunAddr != "" {
		FlagServerAddress = envRunAddr
	}
	if envReportInterval := AppConfig.ReportInterval; envReportInterval != 0 {
		FlagReportInterval = envReportInterval
	}
	if envPollInterval := AppConfig.PollInterval; envPollInterval != 0 {
		FlagPollInterval = envPollInterval
	}
}
