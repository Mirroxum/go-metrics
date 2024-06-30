package main

import (
	"fmt"
	"net/http"
)

func main() {
	loadConfig()

	parseFlags()

	fmt.Printf("HTTP-server address: %s\n", FlagServerAddress)

	err := http.ListenAndServe(FlagServerAddress, MetricRouter())
	if err != nil {
		panic(err)
	}
}
