package fileprocessing

import (
	"bufio"
	"fmt"
	"math"
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/shirou/gopsutil/v4/mem"
)

// FileIndexSummary holds the index map and metadata about the indexed file.
type FileIndexSummary struct {
	Index         map[int]int64
	IndexOffset   int
	NumberOfLines int
}

// memoryLimitFactor defines the fraction of total system memory allowed for index usage.
const memoryLimitFactor = 0.7

// Estimate the memory usage per index entry (approximate).
const bytesPerIndexEntry = 16 // 8 bytes for int key + 8 bytes for int64 value (map overhead excluded)

// GenerateIndex reads a file and generates an index of line numbers
// and their corresponding byte offsets.
// The index is kept in memory and is used to quickly access lines in the file.
// The function takes the file path as argument and maxIndexes to limit the number of indexes.
// If maxIndexes is 0, it calculates the number of indexes that can be generated based on the available memory.
func GenerateIndex(logger *zerolog.Logger, filePath string, maxIndexes int) (*FileIndexSummary, error) {
	// Validate arguments
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	if filePath == "" {
		return nil, errors.New("file path cannot be empty")
	}
	// Open the file
	file, err := os.Open(filePath)
	// Close the file when done
	defer func() {
		if err := file.Close(); err != nil {
			logger.Error().Err(err).Msg("failed to close file")
		}
	}()
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	// Count the number of lines in the file
	linesCount, err := countLines(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to count lines")
	}
	if linesCount == 0 {
		return nil, nil
	}
	// Seek to the beginning of the file
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to seek to beginning of file")
	}

	// If maxIndexes is not provided, calculate the maximum number of indexes
	if maxIndexes == 0 {
		// Determine available memory for index creation
		vmStat, err := mem.VirtualMemory()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get system info")
		}

		// Calculate the maximum number of indexes based on available memory
		availableMemory := float64(vmStat.Available) * memoryLimitFactor
		maxIndexes = int(availableMemory / bytesPerIndexEntry)

		logger.Info().
			Str("memory limit factor", fmt.Sprintf("%.2f", memoryLimitFactor)).
			Str("available memory", fmt.Sprintf("%.2f (GB)", float64(vmStat.Available)/1e9)).
			Str("available memory for index generation", fmt.Sprintf("%.2f (GB)", availableMemory/1e9)).
			Int("maximum number of indexes", maxIndexes).
			Msg("Memory and index statistics")
	}
	if maxIndexes <= 0 {
		return nil, errors.New("insufficient memory available for indexing")
	}
	// Calculate the index offset
	indexOffset := int(math.Ceil(float64(linesCount) / float64(maxIndexes)))

	fileIndexSummary := &FileIndexSummary{
		IndexOffset:   indexOffset,
		NumberOfLines: linesCount,
	}
	indexMap := make(map[int]int64)
	var offset int64 = 0
	currentLine := 0
	scanner := bufio.NewScanner(file)
	// Read the file line by line and populate the index map
	for scanner.Scan() {
		if currentLine%indexOffset == 0 {
			indexMap[currentLine] = offset
		}
		offset += int64(len(scanner.Bytes()) + 1) // +1 for '\n'
		currentLine++
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.Wrap(err, "error reading file")
	}
	fileIndexSummary.Index = indexMap
	return fileIndexSummary, nil
}

// countLines counts the number of lines in a file.
func countLines(file *os.File) (int, error) {
	lineCount := 0
	reader := bufio.NewReader(file)
	for {
		// Read chunks of data
		_, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		lineCount += 1 // Count each newline
	}
	return lineCount, nil
}
