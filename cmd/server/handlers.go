package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Metrics struct {
	ID    string     `json:"id"`              // имя метрики
	MType MetricType `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64     `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64   `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

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
	w.Write([]byte(strings.Join(metrics, "")))
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
	w.Write([]byte(fmt.Sprintf("%v", value)))
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

func UpdateJSONMetricHandler(w http.ResponseWriter, r *http.Request) {

	var metric Metrics
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case Gauge:
		if metric.Value == nil {
			http.Error(w, "Value is required for gauge metric", http.StatusBadRequest)
			return
		}
		storage.UpdateGauge(metric.ID, *metric.Value)
	case Counter:
		if metric.Delta == nil {
			http.Error(w, "Delta is required for counter metric", http.StatusBadRequest)
			return
		}
		storage.UpdateCounter(metric.ID, *metric.Delta)
	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metric); err != nil {
		http.Error(w, "Failed to encode metrics", http.StatusInternalServerError)
		return
	}
}

func GetJSONMetricHandler(w http.ResponseWriter, r *http.Request) {
	var requestMetric Metrics
	if err := json.NewDecoder(r.Body).Decode(&requestMetric); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestMetric.ID == "" || requestMetric.MType == "" {
		http.Error(w, "ID and MType are required", http.StatusBadRequest)
		return
	}
	var responseMetric Metrics
	switch requestMetric.MType {
	case Gauge:
		if value, exists := storage.GetGauge(requestMetric.ID); exists {
			responseMetric = Metrics{
				ID:    requestMetric.ID,
				MType: Gauge,
				Value: &value,
			}
		} else {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}
	case Counter:
		if value, exists := storage.GetCounter(requestMetric.ID); exists {
			responseMetric = Metrics{
				ID:    requestMetric.ID,
				MType: Counter,
				Delta: &value,
			}
		} else {
			http.Error(w, "Metric not found", http.StatusNotFound)
			return
		}
	default:
		http.Error(w, "Invalid metric type", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseMetric); err != nil {
		http.Error(w, "Failed to encode metric", http.StatusInternalServerError)
		return
	}
}
