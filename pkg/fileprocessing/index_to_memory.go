package fileprocessing

import (
	"bufio"
	"fmt"
	"math"
	"os"

	"github.com/pkg/errors"
)

// FileIndexSummary holds the index map and metadata about the indexed file.
type FileIndexSummary struct {
	Index         map[int]int64
	IndexOffset   int
	NumberOfLines int
}

// GenerateIndex reads a file and generates an index of line numbers
// and their corresponding byte offsets.
// The index is kept in memory and is used to quickly access lines in the file.
// The function takes the file path and the maximum number of indexes to generate.
func GenerateIndex(filePath string, maxIndexes int) (*FileIndexSummary, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file %s", filePath)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()
	// Count the number of lines in the file
	linesCount, err := countLines(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed to count lines")
	}
	// No index is generated if file does not have lines or maxIndexes is less than or equal to 0
	if linesCount == 0 || maxIndexes <= 0 {
		return nil, nil
	}
	// Reset the file pointer to the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to seek to beginning of file")
	}
	// Calculate the number of indexes to generate
	indexOffset := math.Ceil(float64(linesCount) / float64(maxIndexes))

	fileIndexSummary := &FileIndexSummary{
		IndexOffset:   int(indexOffset),
		NumberOfLines: linesCount,
	}
	indexMap := map[int]int64{}
	var offset int64 = 0
	currentLine := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Only add the index if the line number is a multiple of indexOffset
		if currentLine%int(indexOffset) == 0 {
			indexMap[currentLine] = offset
		}
		// Increment the offset by the length of the line + 1 for '\n'
		offset += int64(len(scanner.Bytes()) + 1)
		currentLine += 1
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
