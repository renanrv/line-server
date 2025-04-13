//go:build unit

package middlewares_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/renanrv/line-server/pkg/middlewares"
	"github.com/rs/zerolog"
)

func TestLoggingMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectedLog    string
	}{
		{
			name:           "GET request to metrics endpoint",
			method:         "GET",
			url:            "/metrics",
			expectedStatus: http.StatusOK,
			expectedLog:    "",
		},
		{
			name:           "GET request to lines endpoint",
			method:         "GET",
			url:            "/v0/lines/1",
			expectedStatus: http.StatusOK,
			expectedLog:    "endpoint call",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rw := httptest.NewRecorder()

			var logBuffer bytes.Buffer
			log := zerolog.New(&logBuffer)

			middleware := middlewares.LoggingMiddleware(&log)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware(handler).ServeHTTP(rw, req)

			if status := rw.Result().StatusCode; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			logOutput := logBuffer.String()
			if !strings.Contains(logOutput, tt.expectedLog) {
				t.Errorf("expected log output to contain %q, got %q", tt.expectedLog, logOutput)
			}
		})
	}
}
