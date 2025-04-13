//go:build unit

package utils_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/renanrv/line-server/pkg/utils"
)

func TestResponseWriter(t *testing.T) {
	tests := []struct {
		name            string
		statusCode      int
		body            string
		expectedStatus  int
		expectedBytes   int
		expectedHeaders http.Header
	}{
		{
			name:           "Status OK with body",
			statusCode:     http.StatusOK,
			body:           "Hello, world!",
			expectedStatus: http.StatusOK,
			expectedBytes:  13,
			expectedHeaders: http.Header{
				"Content-Type": []string{"text/plain; charset=utf-8"},
			},
		},
		{
			name:            "Status NotFound without body",
			statusCode:      http.StatusNotFound,
			body:            "",
			expectedStatus:  http.StatusNotFound,
			expectedBytes:   0,
			expectedHeaders: http.Header{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			wrappedWriter := utils.WrapResponseWriter(recorder)

			if tt.body != "" {
				recorder.Header().Set("Content-Type", "text/plain; charset=utf-8")
			}

			wrappedWriter.WriteHeader(tt.statusCode)
			if tt.body != "" {
				_, err := wrappedWriter.Write([]byte(tt.body))
				if err != nil {
					t.Errorf("unexpected error while writing body: %v", err)
				}
			}

			if wrappedWriter.Status != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, wrappedWriter.Status)
			}

			if wrappedWriter.BytesWritten != tt.expectedBytes {
				t.Errorf("expected bytes written to be %d, got %d", tt.expectedBytes, wrappedWriter.BytesWritten)
			}

			actualHeaders := wrappedWriter.Header()
			for key, expectedValues := range tt.expectedHeaders {
				if actualValues, ok := actualHeaders[key]; !ok || !reflect.DeepEqual(actualValues, expectedValues) {
					t.Errorf("expected header %s to be %v, got %v", key, expectedValues, actualValues)
				}
			}
		})
	}
}
