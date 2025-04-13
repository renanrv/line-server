//go:build unit

package handler_test

import (
	"context"
	"testing"

	"github.com/renanrv/line-server/services/handler"
	"github.com/renanrv/line-server/services/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	name             string
	request          any
	expectedResponse any
	expectedError    error
}

func TestNew(t *testing.T) {
	logger := zerolog.New(nil)

	_, err := handler.New(&logger)
	if err != nil {
		t.Errorf("New() error = %v, wantErr %v", err, false)
	}
}

func TestHandler_GetV0LinesLineIndex(t *testing.T) {
	tests := []struct {
		name             string
		request          server.GetV0LinesLineIndexRequestObject
		expectedResponse server.GetV0LinesLineIndexResponseObject
		expectedError    error
	}{
		{
			name: "Get existing line",
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: 1,
			},
			expectedResponse: server.GetV0LinesLineIndex200JSONResponse{
				LineResponseJSONResponse: server.LineResponseJSONResponse{
					Text: "text",
				},
			},
			expectedError: nil,
		},
		{
			name: "Invalid line index",
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: -1,
			},
			expectedResponse: server.GetV0LinesLineIndex413JSONResponse{},
			expectedError:    nil,
		},
	}

	logger := zerolog.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, _ := handler.New(&logger)
			ctx := context.Background()
			response, err := h.GetV0LinesLineIndex(ctx, tt.request)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResponse, response)
		})
	}
}
