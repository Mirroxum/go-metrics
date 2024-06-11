package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func UpdateMetricRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/update", func(r chi.Router) {
		r.Post("/{metricType}/{metricName}/{metricValue}", UpdateMetricHandler)
	})
	return r
}

func main() {
	err := http.ListenAndServe(":8080", UpdateMetricRouter())
	if err != nil {
		panic(err)
	}
}
