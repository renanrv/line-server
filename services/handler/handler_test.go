//go:build unit

package handler_test

import (
	"context"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/renanrv/line-server/pkg/fileprocessing"
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
		name             string
		logger           *zerolog.Logger
		filePath         string
		fileIndexSummary *fileprocessing.FileIndexSummary
		expectedError    error
	}{
		{
			name:             "valid dependencies",
			logger:           &zerolog.Logger{},
			filePath:         file.Name(),
			fileIndexSummary: &fileprocessing.FileIndexSummary{},
		},
		{
			name:             "missing logger",
			filePath:         file.Name(),
			fileIndexSummary: &fileprocessing.FileIndexSummary{},
			expectedError:    errors.New("logger is required"),
		},
		{
			name:             "missing file path",
			logger:           &zerolog.Logger{},
			fileIndexSummary: &fileprocessing.FileIndexSummary{},
			expectedError:    errors.New("file path is required"),
		},
		{
			name:             "file does not exist",
			logger:           &zerolog.Logger{},
			filePath:         "test.txt",
			fileIndexSummary: &fileprocessing.FileIndexSummary{},
			expectedError:    errors.New("file does not exist"),
		},
		{
			name:     "missing optional fileIndexSummary",
			logger:   &zerolog.Logger{},
			filePath: file.Name(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.New(tt.logger, tt.filePath, tt.fileIndexSummary)
			if tt.expectedError != nil || err != nil {
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			}
		})
	}
}

type testStruct struct {
	name             string
	fileIndexSummary *fileprocessing.FileIndexSummary
	request          server.GetV0LinesLineIndexRequestObject
	expectedResponse server.GetV0LinesLineIndexResponseObject
	expectedError    error
}

func TestHandler_GetV0LinesLineIndex(t *testing.T) {
	content := "line1\nline2\nline3\n"

	tests := []testStruct{
		{
			name: "Get existing line without file index summary",
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
			name: "Get existing line with indexed line from file index summary",
			fileIndexSummary: &fileprocessing.FileIndexSummary{
				Index: map[int]int64{
					0: 0,
					1: 6,
					2: 12,
				},
				IndexOffset:   1,
				NumberOfLines: 3,
			},
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
			name: "Get existing line without indexed line from file index summary",
			fileIndexSummary: &fileprocessing.FileIndexSummary{
				Index: map[int]int64{
					0: 0,
					2: 12,
				},
				IndexOffset:   2,
				NumberOfLines: 3,
			},
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
			name: "No closest index in file index summary",
			fileIndexSummary: &fileprocessing.FileIndexSummary{
				Index: map[int]int64{
					4: 12,
					5: 18,
				},
				IndexOffset:   2,
				NumberOfLines: 3,
			},
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
			name: "Invalid file path",
			fileIndexSummary: &fileprocessing.FileIndexSummary{
				Index: map[int]int64{
					4: 12,
					5: 18,
				},
				IndexOffset:   1,
				NumberOfLines: 1,
			},
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: 1,
			},
			expectedResponse: nil,
			expectedError:    errors.New("failed to open file"),
		},
		{
			name: "Invalid line index with negative value",
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: -1,
			},
			expectedResponse: server.GetV0LinesLineIndex413Response{},
			expectedError:    nil,
		},
		{
			name: "Invalid line index with negative value with file index summary",
			fileIndexSummary: &fileprocessing.FileIndexSummary{
				Index: map[int]int64{
					0: 0,
					1: 6,
					2: 12,
				},
				IndexOffset:   1,
				NumberOfLines: 3,
			},
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: -1,
			},
			expectedResponse: server.GetV0LinesLineIndex413Response{},
			expectedError:    nil,
		},
		{
			name: "Invalid line index with value greater than number of lines",
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: 3,
			},
			expectedResponse: server.GetV0LinesLineIndex413Response{},
			expectedError:    nil,
		},
		{
			name: "Invalid line index with value greater than number of lines with file index summary",
			fileIndexSummary: &fileprocessing.FileIndexSummary{
				Index: map[int]int64{
					0: 0,
					1: 6,
					2: 12,
				},
				IndexOffset:   1,
				NumberOfLines: 3,
			},
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: 3,
			},
			expectedResponse: server.GetV0LinesLineIndex413Response{},
			expectedError:    nil,
		},
		{
			name: "Invalid index in file index summary",
			fileIndexSummary: &fileprocessing.FileIndexSummary{
				Index:         nil,
				IndexOffset:   2,
				NumberOfLines: 3,
			},
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: 1,
			},
			expectedResponse: nil,
			expectedError:    errors.New("file index is required"),
		},
		{
			name: "Invalid index offset in file index summary",
			fileIndexSummary: &fileprocessing.FileIndexSummary{
				Index: map[int]int64{
					0: 0,
					1: 6,
					2: 12,
				},
				IndexOffset:   0,
				NumberOfLines: 3,
			},
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: 1,
			},
			expectedResponse: nil,
			expectedError:    errors.New("file index offset must be greater than 0"),
		},
		{
			name: "Invalid number of lines in file index summary",
			fileIndexSummary: &fileprocessing.FileIndexSummary{
				Index: map[int]int64{
					0: 0,
					1: 6,
					2: 12,
				},
				IndexOffset:   1,
				NumberOfLines: 0,
			},
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: 1,
			},
			expectedResponse: nil,
			expectedError:    errors.New("file number of lines must be greater than 0"),
		},
		{
			name: "Invalid index position in file index summary",
			fileIndexSummary: &fileprocessing.FileIndexSummary{
				Index: map[int]int64{
					0: 0,
					1: 24,
					2: 12,
				},
				IndexOffset:   1,
				NumberOfLines: 3,
			},
			request: server.GetV0LinesLineIndexRequestObject{
				LineIndex: 1,
			},
			expectedResponse: nil,
			expectedError:    errors.New("error reading file: EOF"),
		},
	}

	logger := zerolog.New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := utils.CreateTempFile(t, content)
			h, err := handler.New(&logger, file.Name(), tt.fileIndexSummary)
			assert.Nil(t, err)
			ctx := context.Background()
			if tt.name == "Invalid file path" {
				os.Remove(file.Name())
			}
			response, err := h.GetV0LinesLineIndex(ctx, tt.request)
			assert.Equal(t, tt.expectedResponse, response)
			if tt.expectedError != nil || err != nil {
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			}
		})
	}
}
