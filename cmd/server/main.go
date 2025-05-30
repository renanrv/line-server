package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/namsral/flag"
	"github.com/pkg/errors"
	"github.com/renanrv/line-server/pkg/fileprocessing"
	"github.com/renanrv/line-server/pkg/middlewares"
	"github.com/renanrv/line-server/services"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

const baseURL = ""

func main() {
	// Define flags
	fs := flag.NewFlagSet("line-server", flag.ExitOnError)

	var (
		debugAddr          = fs.String("debug_addr", ":8081", "debug and metrics listen address")
		httpAddr           = fs.String("http_addr", ":8080", "the address that will expose the server API")
		corsAllowedOrigins = fs.String("cors_allowed_origins", "http://localhost:8080",
			"comma separated list of allowed origins")
		logLevel = fs.Int("log_level", int(zerolog.InfoLevel), "the log level used for logging")
		filePath = fs.String("file_path", "./data/sample_100.txt",
			"the path to the file that will be used to read the lines")
		maxIndexes = fs.Int("max_indexes", 0, "the maximum number of indexes to generate, "+
			"taking into account the limited memory available. If 0, it will use all available memory. If "+
			"negative, it will not generate any indexes.")
	)
	fs.Usage = usageFor(fs, os.Args[0]+" [flags]")
	_ = fs.Parse(os.Args[1:])
	// JSON logger
	zerolog.SetGlobalLevel(zerolog.Level(*logLevel))
	zeroLog := zlog.With().Caller().Str("component", "line-server").Logger()

	// log non-secret arguments to help debugging issues
	zeroLog.Info().
		Str("service", "line-server").
		Str("debug_addr", *debugAddr).
		Str("http_addr", *httpAddr).
		Int("log_level", *logLevel).
		Str("file_path", *filePath).
		Int("max_indexes", *maxIndexes).
		Msg("non-secret arguments")

	zeroLog.Info().Msg("starting line server")

	// Split CORS allowed origins as array of strings
	var corsAllowedOriginsList []string
	if corsAllowedOrigins != nil && *corsAllowedOrigins != "" {
		corsAllowedOriginsList = strings.Split(*corsAllowedOrigins, ",")
	}

	// Add CORS support (Cross Origin Resource Sharing)
	if len(corsAllowedOriginsList) == 0 {
		zeroLog.Warn().
			Err(errors.New("cors_allowed_origins config is empty which defaults to *")).
			Msg("CORS is disabled")
	}

	// Check if indexes should be generated
	var fileIndexSummary *fileprocessing.FileIndexSummary = nil
	if *maxIndexes >= 0 {
		fileIndexSummary, err := fileprocessing.GenerateIndex(&zeroLog, *filePath, *maxIndexes)
		// Validate file index summary
		if err != nil {
			zeroLog.Fatal().Err(err).Msg("failed to generate index")
		}
		zeroLog.Info().Int("length", len(fileIndexSummary.Index)).Msg("index generated successfully")
	}

	dependencies := services.Dependencies{
		Logger:           &zeroLog,
		FilePath:         *filePath,
		FileIndexSummary: fileIndexSummary,
	}
	srv, err := services.New(dependencies)
	if err != nil {
		zeroLog.Fatal().Err(err).Msg("failed to create line-server service")
	}

	mux, err := srv.Router(services.RouterOpts{
		PathPrefix: baseURL,
	})
	if err != nil {
		zeroLog.Fatal().Err(err).Msg("failed to create line-server router")
	}
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: corsAllowedOriginsList,
		AllowedHeaders: []string{"authorization"},
	})
	handlerHTTP := corsHandler.Handler(mux)

	s := &http.Server{
		Addr:    *httpAddr,
		Handler: middlewares.LoggingMiddleware(&zeroLog)(handlerHTTP),
	}

	// Signal handling
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	// Start the application server
	go func() {
		zeroLog.Info().Msgf("starting server on port %s", *httpAddr)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Log error if the server fails to start or if it shuts down unexpectedly
			zeroLog.Error().Err(err).Msg("http server stopped")
		}
	}()

	// Wait for interrupt signal
	<-sigChannel
	zeroLog.Info().Msg("shutting down server")

	err = s.Shutdown(context.Background())
	if err != nil {
		zeroLog.Fatal().Err(err).Msg("failed to gracefully stop the server")
	}
	zeroLog.Info().Msg("server was gracefully stopped")
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			_, err := fmt.Fprintf(w, "\t-%s %s\t%s\n", f.Name, f.DefValue, f.Usage)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error writing flag usage: %v\n", err)
			}
		})
		err := w.Flush()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error flushing tabwriter: %v\n", err)
		}
		fmt.Fprintf(os.Stderr, "\n")
	}
}
