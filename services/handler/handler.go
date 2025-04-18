package handler

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/renanrv/line-server/pkg/fileprocessing"
	"github.com/renanrv/line-server/pkg/utils"
	"github.com/renanrv/line-server/services/server"
	"github.com/rs/zerolog"
)

type Handler struct {
	Logger           *zerolog.Logger
	FilePath         string
	FileIndexSummary *fileprocessing.FileIndexSummary
}

// New function instantiates a handler, checking if all dependencies are valid
func New(l *zerolog.Logger, filePath string, fileIndexSummary *fileprocessing.FileIndexSummary) (server.StrictServerInterface, error) {
	// Validates the handler's dependencies.
	// File index summary is optional, but if provided,
	// it is validated by validateFileIndexSummary function before being used.
	err := validate(l, filePath)
	if err != nil {
		return nil, err
	}
	return Handler{
		Logger:           l,
		FilePath:         filePath,
		FileIndexSummary: fileIndexSummary,
	}, nil
}

// GetV0LinesLineIndex returns a line for a given line index
func (h Handler) GetV0LinesLineIndex(_ context.Context, request server.GetV0LinesLineIndexRequestObject,
) (server.GetV0LinesLineIndexResponseObject, error) {
	// Obtain the result from the file according the requested line index
	text, err := h.readLine(request.LineIndex)
	if err != nil {
		return nil, err
	}
	// Check requested line index and result to infer if invalid index was requested
	if request.LineIndex < 0 || (text == "" && err == nil) {
		return server.GetV0LinesLineIndex413JSONResponse{}, nil
	}
	// Returns successful response
	return server.GetV0LinesLineIndex200JSONResponse{
		LineResponseJSONResponse: server.LineResponseJSONResponse{
			Text: text,
		},
	}, nil
}

// readLine method reads the file and returns the line according the provided line index
func (h Handler) readLine(lineIndex int) (string, error) {
	file, err := os.Open(h.FilePath)
	if err != nil {
		return "", errors.New("could not open file")
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
	}()
	// If no file index summary is available, read the file line by line
	if h.FileIndexSummary == nil {
		currentLine := 0
		return scanFile(lineIndex, file, currentLine)
	}
	// If file index summary is available, seek the line index in the index map
	return h.seekFileLine(file, lineIndex)
}

// seekFileLine function seeks the line index in the file by using the file index map.
// It uses the file index map to find the starting position of the line in the file.
// If the line index is not found in the map, it finds the closest indexed line and seeks to that position.
// It then reads line by line from that position in the file and returns the line when the request index is found.
func (h Handler) seekFileLine(file *os.File, lineIndex int) (string, error) {
	// Validate the file index summary
	err := h.validateFileIndexSummary()
	if err != nil {
		return "", err
	}
	start := int64(0)
	start, ok := h.FileIndexSummary.Index[lineIndex]
	if !ok {
		h.Logger.Info().Int("index", lineIndex).Msg("no index available in index map")
		// Retrieve the closest index for the line index value,
		// in order to identify the best indexed starting point to seek the line
		currentLine, err := h.closestIndexedFileLine(lineIndex)
		if err != nil {
			return "", err
		}
		h.Logger.Debug().Int("index", lineIndex).Int("closest_index", currentLine).Msg("closest index available in index map")
		start, ok = h.FileIndexSummary.Index[currentLine]
		if !ok {
			h.Logger.Warn().Int("index", lineIndex).Msg("no closest index available in index map")
			// If no closest index is available, read the file line by line from the beginning
			_, err = file.Seek(0, 0)
			if err != nil {
				return "", errors.Wrap(err, "failed to seek to beginning of file")
			}
			return scanFile(lineIndex, file, currentLine)
		}
		h.Logger.Debug().Int("index", lineIndex).Int64("start", start).Msg("Position in file")
		_, err = file.Seek(start, 0)
		if err != nil {
			return "", errors.Wrap(err, "failed to seek to index position")
		}
		return scanFile(lineIndex, file, currentLine)
	}
	// If the line index is found in the index map, seek to that position
	_, err = file.Seek(start, 0)
	if err != nil {
		return "", errors.Wrap(err, "failed to seek to index position")
	}
	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", errors.Wrap(err, "error reading file")
	}
	// Return the line value without the `\n` character
	if len(line) > 0 && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}
	return line, nil
}

// closestIndexedFileLine fetches the closest index for the requested line index
func (h Handler) closestIndexedFileLine(lineIndex int) (int, error) {
	// FileIndexSummary has been validated before,
	// so we can safely use it here
	numberOfLines := h.FileIndexSummary.NumberOfLines
	indexOffset := h.FileIndexSummary.IndexOffset
	for indexedFileLine := 0; indexedFileLine < numberOfLines; indexedFileLine += indexOffset {
		if (indexedFileLine <= lineIndex) && (indexedFileLine+indexOffset > lineIndex) {
			return indexedFileLine, nil
		}
	}
	return 0, errors.New("no file index found")
}

// validateFileIndexSummary validates the file index summary
func (h Handler) validateFileIndexSummary() error {
	if h.FileIndexSummary == nil {
		return errors.New("file index summary is required")
	}
	if h.FileIndexSummary.Index == nil {
		return errors.New("file index is required")
	}
	if h.FileIndexSummary.IndexOffset <= 0 {
		return errors.New("file index offset must be greater than 0")
	}
	if h.FileIndexSummary.NumberOfLines <= 0 {
		return errors.New("file number of lines must be greater than 0")
	}
	return nil
}

// scanFile function uses a scanner to read a file line by line and return the line from the provided line index
func scanFile(lineIndex int, file *os.File, currentLine int) (string, error) {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if currentLine == lineIndex {
			return scanner.Text(), nil
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return "", errors.Wrap(err, "Error reading file")
	}

	// Line index out of range
	return "", nil
}

// validate function validates the dependencies to instantiate a new handler
func validate(logger *zerolog.Logger, filePath string) error {
	if logger == nil {
		return errors.New("logger is required")
	}
	if filePath == "" {
		return errors.New("file path is required")
	}
	if !utils.FileExists(filePath) {
		return errors.New("file does not exist")
	}
	return nil
}
