package main

import "flag"

var (
	FlagServerAddress  = flag.String("a", "localhost:8080", "HTTP-server address")
	FlagReportInterval = flag.Int64("r", 10, "Frequency of sending metrics to the server")
	FlagPollInterval   = flag.Int64("p", 2, "Frequency of polling metrics from the runtime package")
)
