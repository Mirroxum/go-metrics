package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

type RuntimeMetrics struct {
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64
	PollCount     int64
	RandomValue   float64
}

func updateRuntimeMetrics() RuntimeMetrics {
	var metrics RuntimeMetrics
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	metrics.Alloc = float64(stats.Alloc)
	metrics.BuckHashSys = float64(stats.BuckHashSys)
	metrics.Frees = float64(stats.Frees)
	metrics.GCCPUFraction = stats.GCCPUFraction
	metrics.GCSys = float64(stats.GCSys)
	metrics.HeapAlloc = float64(stats.HeapAlloc)
	metrics.HeapIdle = float64(stats.HeapIdle)
	metrics.HeapInuse = float64(stats.HeapInuse)
	metrics.HeapObjects = float64(stats.HeapObjects)
	metrics.HeapReleased = float64(stats.HeapReleased)
	metrics.HeapSys = float64(stats.HeapSys)
	metrics.LastGC = float64(stats.LastGC)
	metrics.Lookups = float64(stats.Lookups)
	metrics.MCacheInuse = float64(stats.MCacheInuse)
	metrics.MCacheSys = float64(stats.MCacheSys)
	metrics.MSpanInuse = float64(stats.MSpanInuse)
	metrics.MSpanSys = float64(stats.MSpanSys)
	metrics.Mallocs = float64(stats.Mallocs)
	metrics.NextGC = float64(stats.NextGC)
	metrics.NumForcedGC = float64(stats.NumForcedGC)
	metrics.NumGC = float64(stats.NumGC)
	metrics.OtherSys = float64(stats.OtherSys)
	metrics.PauseTotalNs = float64(stats.PauseTotalNs)
	metrics.StackInuse = float64(stats.StackInuse)
	metrics.StackSys = float64(stats.StackSys)
	metrics.Sys = float64(stats.Sys)
	metrics.TotalAlloc = float64(stats.TotalAlloc)
	metrics.PollCount = 1
	metrics.RandomValue = rand.Float64() * 100
	return metrics
}

func sendDataToServer(serverURL string, metrics RuntimeMetrics) error {
	metricsMap := map[string]interface{}{
		"Alloc":         metrics.Alloc,
		"BuckHashSys":   metrics.BuckHashSys,
		"Frees":         metrics.Frees,
		"GCCPUFraction": metrics.GCCPUFraction,
		"GCSys":         metrics.GCSys,
		"HeapAlloc":     metrics.HeapAlloc,
		"HeapIdle":      metrics.HeapIdle,
		"HeapInuse":     metrics.HeapInuse,
		"HeapObjects":   metrics.HeapObjects,
		"HeapReleased":  metrics.HeapReleased,
		"HeapSys":       metrics.HeapSys,
		"LastGC":        metrics.LastGC,
		"Lookups":       metrics.Lookups,
		"MCacheInuse":   metrics.MCacheInuse,
		"MCacheSys":     metrics.MCacheSys,
		"MSpanInuse":    metrics.MSpanInuse,
		"MSpanSys":      metrics.MSpanSys,
		"Mallocs":       metrics.Mallocs,
		"NextGC":        metrics.NextGC,
		"NumForcedGC":   metrics.NumForcedGC,
		"NumGC":         metrics.NumGC,
		"OtherSys":      metrics.OtherSys,
		"PauseTotalNs":  metrics.PauseTotalNs,
		"StackInuse":    metrics.StackInuse,
		"StackSys":      metrics.StackSys,
		"Sys":           metrics.Sys,
		"TotalAlloc":    metrics.TotalAlloc,
		"PollCount":     metrics.PollCount,
		"RandomValue":   metrics.RandomValue,
	}
	type Metrics struct {
		ID    string   `json:"id"`              // имя метрики
		MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
		Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
		Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
	}

	for metricName, metricValue := range metricsMap {
		var metricType string
		if metricName == "PollCount" {
			metricType = "counter"
		} else {
			metricType = "gauge"
		}

		metric := Metrics{
			ID:    metricName,
			MType: metricType,
		}

		if metricType == "gauge" {
			value := metricValue.(float64)
			metric.Value = &value
		} else {
			delta := metricValue.(int64)
			metric.Delta = &delta
		}

		jsonData, err := json.Marshal(metric)
		if err != nil {
			return fmt.Errorf("failed to marshal metric: %w", err)
		}

		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		if _, err := gz.Write(jsonData); err != nil {
			return fmt.Errorf("failed to compress data: %w", err)
		}
		gz.Close()

		link := fmt.Sprintf("%s/update/", serverURL)
		fmt.Println("Sending to:", link, "Data:", string(jsonData))
		req, err := http.NewRequest("POST", link, &buf)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept-Encoding", "gzip")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			return fmt.Errorf("failed to send metric: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("server returned non-200 status code: %d", resp.StatusCode)
		}
	}

	return nil
}

func main() {
	loadConfig()
	parseFlags()

	var (
		pollInterval   = time.Duration(FlagPollInterval) * time.Second
		reportInterval = time.Duration(FlagReportInterval) * time.Second
		serverAddress  = fmt.Sprintf("http://%s", FlagServerAddress)
	)

	fmt.Printf("Metrics will be sent to the HTTP server address: %s\n", serverAddress)
	fmt.Printf("Frequency of sending metrics to the server: %d seconds\n", FlagReportInterval)
	fmt.Printf("Frequency of polling metrics from the runtime package: %d seconds\n", FlagPollInterval)

	lastReportTime := time.Now()

	for {
		time.Sleep(pollInterval)
		metrics := updateRuntimeMetrics()

		if time.Since(lastReportTime) >= reportInterval {
			err := sendDataToServer(serverAddress, metrics)
			if err != nil {
				fmt.Println("Error sending metrics to server:", err)
			} else {
				fmt.Println("Metrics sent to server.")
			}
			lastReportTime = time.Now()
		}
	}
}
