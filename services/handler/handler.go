package handler

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/renanrv/line-server/pkg/utils"
	"github.com/renanrv/line-server/services/server"
	"github.com/rs/zerolog"
)

type Handler struct {
	Logger   *zerolog.Logger
	FilePath string
}

// New function instantiates a handler, checking if all dependencies are valid
func New(l *zerolog.Logger, filePath string) (server.StrictServerInterface, error) {
	// Validates the handler's dependencies
	err := validate(l, filePath)
	if err != nil {
		return nil, err
	}
	return Handler{
		Logger:   l,
		FilePath: filePath,
	}, nil
}

// GetV0LinesLineIndex returns a line for a given line index
func (h Handler) GetV0LinesLineIndex(ctx context.Context, request server.GetV0LinesLineIndexRequestObject,
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
		return "", errors.New("Could not open file")
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}()
	return scanFile(lineIndex, file)
}

// scanFile function uses a scanner to read a file line by line and return the line from the provided line index
func scanFile(lineIndex int, file *os.File) (string, error) {
	scanner := bufio.NewScanner(file)
	current := 0
	for scanner.Scan() {
		if current == lineIndex {
			return scanner.Text(), nil
		}
		current++
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
