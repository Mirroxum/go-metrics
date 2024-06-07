package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func UpdateMetricHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method + " " + r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 5 {
		http.Error(w, "invalid URL format", http.StatusNotFound)
		return
	}

	metricType := MetricType(parts[2])
	metricName := parts[3]
	metricValue := parts[4]

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
