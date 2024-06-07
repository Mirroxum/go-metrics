package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateMetricHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{"GaugeOK", "POST", "/update/gauge/testGauge/12.34", http.StatusOK},
		{"CounterOK", "POST", "/update/counter/testCounter/10", http.StatusOK},
		{"GetNotAllowed", "GET", "/update/gauge/testGauge/12.34", http.StatusMethodNotAllowed},
		{"invalidType", "POST", "/update/invalidType/testGauge/12.34", http.StatusBadRequest},
		{"invalidValue", "POST", "/update/gauge/testGauge/invalidValue", http.StatusBadRequest},
		{"emptyType", "POST", "/update/gauge//12.34", http.StatusNotFound},
		{"emptyValue", "POST", "/update/gauge/testGauge", http.StatusNotFound},
		{"extraURL", "POST", "/update/gauge/testGauge/12.34/extra", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			handler := http.HandlerFunc(UpdateMetricHandler)
			handler(w, req)

			assert.Equal(t, w.Code, tt.expectedStatus)
		})
	}
}
