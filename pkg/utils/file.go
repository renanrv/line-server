package utils

import "os"

// FileExists function checks the file's metadata and if the file does not exist in case of error
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
