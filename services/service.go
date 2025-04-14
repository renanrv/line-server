//go:generate go tool oapi-codegen --config=docs/openapi/oapi-codegen.config.yaml docs/openapi/lineserver.openapi.yaml

package services

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/renanrv/line-server/services/handler"
	"github.com/renanrv/line-server/services/server"
	"github.com/rs/zerolog"
)

type Dependencies struct {
	Logger   *zerolog.Logger
	FilePath string
}

type service struct {
	logger   *zerolog.Logger
	filePath string
}

// RouterOpts represents router options
type RouterOpts struct {
	PathPrefix     string
	ExistingRouter *http.ServeMux // Optional
}

// Service is the Service Interface
type Service interface {
	Router(opts RouterOpts) (*http.ServeMux, error)
}

// New creates a new service
func New(d Dependencies) (Service, error) {
	if d.Logger == nil {
		return nil, errors.New("logger is required")
	}
	if d.FilePath == "" {
		return nil, errors.New("file path is required")
	}
	return service{
		logger:   d.Logger,
		filePath: d.FilePath,
	}, nil
}

// Router returns a router configured with the quantifier service
// router param is optional, if nil a new router is created
// if a router is passed the new Router gets augmented with the quantifier service
func (s service) Router(opts RouterOpts) (*http.ServeMux, error) {
	h, err := handler.New(s.logger, s.filePath)
	if err != nil {
		return nil, err
	}

	router := http.NewServeMux()

	handlerOptions := server.StdHTTPServerOptions{
		BaseURL:    opts.PathPrefix,
		BaseRouter: router,
	}
	hdl := server.NewStrictHandler(h, nil)
	server.HandlerWithOptions(hdl, handlerOptions)

	return router, nil
}
