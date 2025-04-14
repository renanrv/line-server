package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

const filePath = "./internal-tools/file-generator/output/output.txt"

// FileGenerator generates a file with a specified number of lines to be used for performance tests
func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zeroLog := zlog.With().Caller().Str("component", "file-generator").Logger()

	if len(os.Args) < 2 {
		zeroLog.Fatal().Msg("Usage: go run internal-tools/file-generator/main.go <number_of_lines>")
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil || n < 0 {
		zeroLog.Fatal().Str("number_of_lines", os.Args[1]).Msg("Invalid number of lines")
	}

	file, err := os.Create(filePath)
	if err != nil {
		zeroLog.Fatal().Err(err).Msg("Could not create output file")
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()

	for i := 0; i < n; i++ {
		line := fmt.Sprintf("Line %d\n", i)
		_, err := file.WriteString(line)
		if err != nil {
			zeroLog.Fatal().Err(err).Msg("Could not write to output file")
		}
	}
	zeroLog.Info().Msg(fmt.Sprintf("Successfully wrote %d lines to %s\n", n, filePath))
}
