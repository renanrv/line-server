//go:build unit

package fileprocessing_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/renanrv/line-server/pkg/fileprocessing"
	"github.com/renanrv/line-server/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateIndex(t *testing.T) {
	tests := []struct {
		name                     string
		content                  string
		maxIndexes               int
		expectedFileIndexSummary *fileprocessing.FileIndexSummary
		expectedError            error
	}{
		{
			name:                     "Empty file",
			content:                  "",
			maxIndexes:               10,
			expectedFileIndexSummary: nil,
			expectedError:            nil,
		},
		{
			name:       "Single line file",
			content:    "line1\n",
			maxIndexes: 10,
			expectedFileIndexSummary: &fileprocessing.FileIndexSummary{
				Index: map[int]int64{
					0: 0,
				},
				IndexOffset:   1,
				NumberOfLines: 1,
			},
			expectedError: nil,
		},
		{
			name:       "Multiple lines with maxIndexes > lines",
			content:    "line1\nline2\nline3\n",
			maxIndexes: 10,
			expectedFileIndexSummary: &fileprocessing.FileIndexSummary{
				Index: map[int]int64{
					0: 0,
					1: 6,
					2: 12,
				},
				IndexOffset:   1,
				NumberOfLines: 3,
			},
			expectedError: nil,
		},
		{
			name:       "Multiple lines with maxIndexes < lines",
			content:    "line1\nline2\nline3\nline4\nline5\n",
			maxIndexes: 2,
			expectedFileIndexSummary: &fileprocessing.FileIndexSummary{
				Index: map[int]int64{
					0: 0,
					3: 18,
				},
				IndexOffset:   3,
				NumberOfLines: 5,
			},
			expectedError: nil,
		},
		{
			name:                     "File does not exist",
			content:                  "",
			maxIndexes:               10,
			expectedFileIndexSummary: nil,
			expectedError: errors.New("failed to open file nonexistent_file.txt: " +
				"open nonexistent_file.txt: no such file or directory"),
		},
		{
			name:                     "maxIndexes is 0",
			content:                  "line1\nline2\nline3\n",
			maxIndexes:               0,
			expectedFileIndexSummary: nil,
			expectedError:            nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string
			if tt.name != "File does not exist" {
				file := utils.CreateTempFile(t, tt.content)
				filePath = file.Name()
			} else {
				filePath = "nonexistent_file.txt"
			}

			result, err := fileprocessing.GenerateIndex(filePath, tt.maxIndexes)

			assert.Equal(t, tt.expectedFileIndexSummary, result)
			if tt.expectedError != nil || err != nil {
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			}
		})
	}
}
