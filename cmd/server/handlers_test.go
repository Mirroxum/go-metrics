package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetricRouter(t *testing.T) {
	ts := httptest.NewServer(MetricRouter())
	defer ts.Close()

	tests := []struct {
		name           string
		method         string
		url            string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{"GaugeOK", "POST", "/update/gauge/testGauge/12.34", "", http.StatusOK, ""},
		{"CounterOK", "POST", "/update/counter/testCounter/10", "", http.StatusOK, ""},
		{"GetNotAllowed", "GET", "/update/gauge/testGauge/12.34", "", http.StatusMethodNotAllowed, ""},
		{"invalidType", "POST", "/update/invalidType/testGauge/12.34", "", http.StatusBadRequest, ""},
		{"invalidValue", "POST", "/update/gauge/testGauge/invalidValue", "", http.StatusBadRequest, ""},
		{"emptyType", "POST", "/update/gauge//12.34", "", http.StatusNotFound, ""},
		{"emptyValue", "POST", "/update/gauge/testGauge", "", http.StatusNotFound, ""},
		{"extraURL", "POST", "/update/gauge/testGauge/12.34/extra", "", http.StatusNotFound, ""},

		{"PostGaugeOK", "POST", "/update", `{"id":"TestGauge", "type":"gauge", "value":12.34}`, http.StatusOK, `{"id":"TestGauge", "type":"gauge", "value":12.34}`},
		{"PostCounterOK", "POST", "/update", `{"id":"TestCounter", "type":"counter", "delta":10}`, http.StatusOK, `{"id":"TestCounter", "type":"counter", "delta":10}`},
		{"PostInvalidType", "POST", "/update", `{"id":"TestGauge", "type":"invalidType", "value":12.34}`, http.StatusBadRequest, ""},
		{"PostMissingValue", "POST", "/update", `{"id":"TestGauge", "type":"gauge"}`, http.StatusBadRequest, ""},
		{"PostMissingDelta", "POST", "/update", `{"id":"TestCounter", "type":"counter"}`, http.StatusBadRequest, ""},
		{"PostInvalidJSON", "POST", "/update", `invalid json`, http.StatusBadRequest, ""},

		{"GetGaugeOK", "GET", "/value/gauge/testGauge", "", http.StatusOK, "12.34"},
		{"GetCounterOK", "GET", "/value/counter/testCounter", "", http.StatusOK, "10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tt.method == "POST" && tt.url == "/update" {
				req, err = http.NewRequest(tt.method, ts.URL+tt.url, bytes.NewBufferString(tt.requestBody))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req, err = http.NewRequest(tt.method, ts.URL+tt.url, nil)
			}
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.method == "GET" || (tt.method == "POST" && tt.url == "/update") {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err)

				if tt.expectedBody != "" {
					assert.JSONEq(t, tt.expectedBody, string(body))
				}
			}
		})
	}
}
