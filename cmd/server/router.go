package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func MetricRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(DecompressGzipMiddleware)
	r.Get("/", GetMetricsHandler)
	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{metricValue}", UpdateMetricHandler)
		r.Post("/", UpdateJSONMetricHandler)
	})
	r.Route("/value", func(r chi.Router) {
		r.Get("/{metricType}/{metricName}", GetMetricHandler)
		r.Post("/", GetJSONMetricHandler)
	})
	return r
}
