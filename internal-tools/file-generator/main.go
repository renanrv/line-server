package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

const (
	BytesInGB = 1024 * 1024 * 1024
)

// main generates a file of a specified size in GB with lines numbered from 0 to n.
func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zeroLog := zlog.With().Caller().Str("component", "file-generator").Logger()

	if len(os.Args) < 3 {
		zeroLog.Fatal().Msg("Usage: go run internal-tools/file-generator/main.go <file_size_gb> <output_file_path>")
	}

	sizeGB, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil || sizeGB <= 0 {
		zeroLog.Fatal().Str("file_size_gb", os.Args[1]).Msg("Invalid file size in GB")
	}
	targetSize := int64(sizeGB * BytesInGB)

	filePath := os.Args[2]
	if filePath == "" {
		zeroLog.Fatal().Str("output_file_path", filePath).Msg("Invalid output file path")
	}

	file, err := os.Create(filePath)
	if err != nil {
		zeroLog.Fatal().Err(err).Msg("Could not create output file")
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v", err)
		}
	}()

	var totalBytesWritten int64
	var lineCount int
	for totalBytesWritten < targetSize {
		line := fmt.Sprintf("Line %d\n", lineCount)
		n, err := file.WriteString(line)
		if err != nil {
			zeroLog.Fatal().Err(err).Msg("Could not write to output file")
		}
		totalBytesWritten += int64(n)
		lineCount++
	}

	zeroLog.Info().Msg(
		fmt.Sprintf("Successfully wrote approximately %.2f GB (%d lines) to %s",
			float64(totalBytesWritten)/BytesInGB, lineCount, filePath))
}
