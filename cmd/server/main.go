package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MetricRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", GetMetricsHandler)
	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{metricValue}", UpdateMetricHandler)
	})
	r.Route("/value", func(r chi.Router) {
		r.Get("/{metricType}/{metricName}", GetMetricHandler)
	})
	return r
}

func main() {
	err := http.ListenAndServe(":8080", MetricRouter())
	if err != nil {
		panic(err)
	}
}
