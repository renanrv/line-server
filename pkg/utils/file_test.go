//go:build unit

package utils_test

import (
	"os"
	"testing"

	"github.com/renanrv/line-server/pkg/utils"
)

func TestFileExists(t *testing.T) {
	t.Run("existing file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "testfile.txt")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())
		defer tmpFile.Close()

		if !utils.FileExists(tmpFile.Name()) {
			t.Errorf("Expected file to exist, but FileExists returned false")
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		nonExistentPath := "not_a_real_file.txt"
		if utils.FileExists(nonExistentPath) {
			t.Errorf("Expected file not to exist, but FileExists returned true")
		}
	})
}
