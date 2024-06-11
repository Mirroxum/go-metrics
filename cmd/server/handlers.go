package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

func GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var metrics []string

	gauges := storage.gauges
	for name, value := range gauges {
		metrics = append(metrics, fmt.Sprintf("%s (Gauge): %.2f\n", name, value))
	}

	counters := storage.counters
	for name, value := range counters {
		metrics = append(metrics, fmt.Sprintf("%s (Counter): %d\n", name, value))
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, strings.Join(metrics, ""))
}

func GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	metricType := MetricType(chi.URLParam(r, "metricType"))
	metricName := chi.URLParam(r, "metricName")

	var value interface{}
	var exists bool

	switch metricType {
	case Gauge:
		value, exists = storage.GetGauge(metricName)
	case Counter:
		value, exists = storage.GetCounter(metricName)
	default:
		http.Error(w, "invalid metric type", http.StatusBadRequest)
		return
	}

	if !exists {
		http.Error(w, "metric not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v", value)
}

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
