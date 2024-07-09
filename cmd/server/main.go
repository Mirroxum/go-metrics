package main

import (
	"github.com/Mirroxum/go-metrics/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	loadConfig()
	parseFlags()

	if errLog := logger.Initialize(FlagLogLevel); errLog != nil {
		panic(errLog)
	}

	logger.Log.Info("HTTP-server", zap.String("address", FlagServerAddress))

	err := http.ListenAndServe(FlagServerAddress, MetricRouter())
	if err != nil {
		panic(err)
	}
}
