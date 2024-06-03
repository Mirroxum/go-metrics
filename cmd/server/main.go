package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MetricType string

const (
	Gauge   MetricType = "gauge"
	Counter MetricType = "counter"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (s *MemStorage) UpdateGauge(name string, value float64) {
	s.gauges[name] = value
}

func (s *MemStorage) UpdateCounter(name string, value int64) {
	s.counters[name] += value
}

func (s *MemStorage) GetGauge(name string, value float64) (float64, bool) {
	value, exists := s.gauges[name]
	return value, exists
}

func (s *MemStorage) GetCounter(name string, value int64) (int64, bool) {
	value, exists := s.counters[name]
	return value, exists
}

var storage = NewMemStorage()

func updateMetricHandler(w http.ResponseWriter, r *http.Request) {
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
		}
		storage.UpdateGauge(metricName, value)
	case Counter:
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "invalid counter value", http.StatusBadRequest)
		}
		storage.UpdateCounter(metricName, value)
	default:
		http.Error(w, "invalid metric type", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, updateMetricHandler)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
