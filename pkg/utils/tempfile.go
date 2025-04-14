package utils

import (
	"os"
	"testing"
)

// CreateTempFile creates a temporary file for tests with the given content and returns the file handle.
// It also schedules the file for deletion after the test completes.
func CreateTempFile(t *testing.T, content string) *os.File {
	// Marks test as a helper for better test failure logs
	t.Helper()

	tmpFile, err := os.CreateTemp("", "testfile.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Close the file so it could be read afterward
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Schedule deletion
	t.Cleanup(func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	})

	return tmpFile
}
