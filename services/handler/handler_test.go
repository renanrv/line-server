//go:build unit

package handler_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/renanrv/line-server/pkg/utils"
	"github.com/renanrv/line-server/services/handler"
	"github.com/renanrv/line-server/services/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	content := "line1\nline2\nline3\n"
	file := utils.CreateTempFile(t, content)

	tests := []struct {
		name          string
		logger        *zerolog.Logger
		filePath      string
		expectedError error
	}{
		{
			name:     "valid dependencies",
			logger:   &zerolog.Logger{},
			filePath: file.Name(),
		},
		{
			name:          "missing logger",
			filePath:      file.Name(),
			expectedError: errors.New("logger is required"),
		},
		{
			name:          "missing file path",
			logger:        &zerolog.Logger{},
			expectedError: errors.New("file path is required"),
		},
		{
			name:          "file does not exist",
			logger:        &zerolog.Logger{},
			filePath:      "test.txt",
			expectedError: errors.New("file does not exist"),
		},
	}
	logger := zerolog.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.New(&logger, file.Name())
			if err != nil {
				t.Errorf("New() error = %v, wantErr %v", err, false)
			}
		})
	}
}

type testStruct struct {
	name             string
	filePath         string
	request          server.GetV0LinesLineIndexRequestObject
	expectedResponse server.GetV0LinesLineIndexResponseObject
	expectedError    error
}

func TestHandler_GetV0LinesLineIndex(t *testing.T) {
	content := "line1\nline2\nline3\n"
	file := utils.CreateTempFile(t, content)

	tests := []testStruct{
		{
			name:     "Get existing line",
			filePath: file.Name(),
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: 1,
			},
			expectedResponse: server.GetV0LinesLineIndex200JSONResponse{
				LineResponseJSONResponse: server.LineResponseJSONResponse{
					Text: "line2",
				},
			},
			expectedError: nil,
		},
		{
			name:     "Invalid line index with negative value",
			filePath: file.Name(),
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: -1,
			},
			expectedResponse: server.GetV0LinesLineIndex413JSONResponse{},
			expectedError:    nil,
		},
		{
			name:     "Invalid line index with value greater than number of lines",
			filePath: file.Name(),
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: 3,
			},
			expectedResponse: server.GetV0LinesLineIndex413JSONResponse{},
			expectedError:    nil,
		},
	}

	logger := zerolog.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, _ := handler.New(&logger, tt.filePath)
			ctx := context.Background()
			response, err := h.GetV0LinesLineIndex(ctx, tt.request)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedResponse, response)
		})
	}
}
