package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := MetricType(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")
	metricValue := chi.URLParam(r, "metricValue")

	if metricName == "" {
		http.Error(w, "metric name is required", http.StatusNotFound)
		return
	}

	switch metricType {
	case Gauge:
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "invalid gauge value", http.StatusBadRequest)
			return
		}
		storage.UpdateGauge(metricName, value)
	case Counter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "invalid counter value", http.StatusBadRequest)
			return
		}
		storage.UpdateCounter(metricName, value)
	default:
		http.Error(w, "invalid metric type", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
