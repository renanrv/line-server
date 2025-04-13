package handler

import (
	"context"

	"github.com/renanrv/line-server/services/server"
	"github.com/rs/zerolog"
)

type Handler struct {
	Logger *zerolog.Logger
}

func New(l *zerolog.Logger) (server.StrictServerInterface, error) {
	return Handler{
		Logger: l,
	}, nil
}

// GetV0LinesLineIndex returns a line for a given line index
func (h Handler) GetV0LinesLineIndex(ctx context.Context, request server.GetV0LinesLineIndexRequestObject,
) (server.GetV0LinesLineIndexResponseObject, error) {
	// TODO: Validate line index
	fileLineCount := 100
	if request.LineIndex < 0 || request.LineIndex >= fileLineCount {
		return server.GetV0LinesLineIndex413JSONResponse{}, nil
	}
	return server.GetV0LinesLineIndex200JSONResponse{
		LineResponseJSONResponse: server.LineResponseJSONResponse{
			Text: "text",
		},
	}, nil
}
