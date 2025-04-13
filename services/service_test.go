//go:build unit

package services_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/renanrv/line-server/services"
	"github.com/rs/zerolog"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name             string
		dependencies     services.Dependencies
		expectedErrorMsg string
	}{
		{
			name: "Valid dependencies",
			dependencies: services.Dependencies{
				Logger: &zerolog.Logger{},
			},
			expectedErrorMsg: "",
		},
		{
			name: "Nil logger",
			dependencies: services.Dependencies{
				Logger: nil,
			},
			expectedErrorMsg: "logger is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := services.New(tt.dependencies)
			if tt.expectedErrorMsg != "" {
				if err == nil {
					t.Errorf("Expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.expectedErrorMsg) {
					t.Errorf("Expected error message to contain \"%s\", got \"%s\"", tt.expectedErrorMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("Did not expect an error, got %v", err)
			}
		})
	}
}

func TestRouter(t *testing.T) {
	tests := []struct {
		name        string
		pathPrefix  string
		expectError bool
	}{
		{
			name:        "Without existing router",
			pathPrefix:  "/api",
			expectError: false,
		},
		{
			name:        "With existing router",
			pathPrefix:  "/api",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zerolog.New(nil)
			svc, err := services.New(services.Dependencies{Logger: &logger})
			if err != nil {
				t.Fatalf("Failed to create svc: %v", err)
			}

			var existingRouter *http.ServeMux
			if tt.name == "With existing router" {
				existingRouter = http.NewServeMux()
			}

			router, err := svc.Router(services.RouterOpts{PathPrefix: tt.pathPrefix, ExistingRouter: existingRouter})
			if (err != nil) != tt.expectError {
				t.Errorf("Router() error = %v, expectError %v", err, tt.expectError)
			}

			if router == nil {
				t.Error("Expected non-nil router")
			}
		})
	}
}
