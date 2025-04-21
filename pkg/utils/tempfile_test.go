//go:build unit

package utils_test

import (
	"os"
	"testing"

	"github.com/renanrv/line-server/pkg/utils"
)

func TestCreateTempFile(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"empty content", ""},
		{"multiple lines", "line1\nline2\nline3"},
	}

	// store temp file names so we can check they are deleted later
	var fileNames []string

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile := utils.CreateTempFile(t, tt.content)
			fileNames = append(fileNames, tmpFile.Name())

			// Check that file exists
			if !utils.FileExists(tmpFile.Name()) {
				t.Errorf("Temp file %s does not exist during test", tmpFile.Name())
			}

			// Check content
			data, err := os.ReadFile(tmpFile.Name())
			if err != nil {
				t.Fatalf("Failed to read temp file: %v", err)
			}
			if string(data) != tt.content {
				t.Errorf("FilePath content mismatch: got %q, want %q", string(data), tt.content)
			}
		})
	}

	// Check after all subtests (and t.Cleanup calls)
	t.Cleanup(func() {
		for _, name := range fileNames {
			if utils.FileExists(name) {
				t.Errorf("Expected temp file %s to be deleted, but it still exists", name)
			}
		}
	})
}
