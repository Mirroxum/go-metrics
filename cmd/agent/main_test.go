package main

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendDataToServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	serverURL := server.URL

	metrics := RuntimeMetrics{
		Alloc:         rand.Float64() * 100,
		BuckHashSys:   rand.Float64() * 100,
		Frees:         rand.Float64() * 100,
		GCCPUFraction: rand.Float64() * 100,
		GCSys:         rand.Float64() * 100,
		HeapAlloc:     rand.Float64() * 100,
		HeapIdle:      rand.Float64() * 100,
		HeapInuse:     rand.Float64() * 100,
		HeapObjects:   rand.Float64() * 100,
		HeapReleased:  rand.Float64() * 100,
		HeapSys:       rand.Float64() * 100,
		LastGC:        rand.Float64() * 100,
		Lookups:       rand.Float64() * 100,
		MCacheInuse:   rand.Float64() * 100,
		MCacheSys:     rand.Float64() * 100,
		MSpanInuse:    rand.Float64() * 100,
		MSpanSys:      rand.Float64() * 100,
		Mallocs:       rand.Float64() * 100,
		NextGC:        rand.Float64() * 100,
		NumForcedGC:   rand.Float64() * 100,
		NumGC:         rand.Float64() * 100,
		OtherSys:      rand.Float64() * 100,
		PauseTotalNs:  rand.Float64() * 100,
		StackInuse:    rand.Float64() * 100,
		StackSys:      rand.Float64() * 100,
		Sys:           rand.Float64() * 100,
		TotalAlloc:    rand.Float64() * 100,
		PollCount:     1,
		RandomValue:   rand.Float64() * 100,
	}

	err := sendDataToServer(serverURL, metrics)
	assert.NoError(t, err)
}

func TestUpdateRuntimeMetrics(t *testing.T) {
	metrics := updateRuntimeMetrics()

	assert.NotZero(t, metrics.Alloc, "Expected non-zero Alloc value")
	assert.NotZero(t, metrics.RandomValue, "Expected non-zero RandomValue value")
}
