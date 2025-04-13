package middlewares

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/renanrv/line-server/pkg/utils"
	"github.com/rs/zerolog"
)

type contextKey string

const RequestTraceIDKey contextKey = "RequestTraceID"

// LoggingMiddleware Logs the status code and the request duration
func LoggingMiddleware(log *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := utils.WrapResponseWriter(w)
			start := time.Now()

			// Add trace id for logging
			requestTraceID := r.Header.Get("x-trace-id")
			// If the header does not exist, inject a new trace-id
			if requestTraceID == "" {
				requestTraceID = uuid.New().String()
			}

			// Inject trace-id into context
			ctx := context.WithValue(r.Context(), RequestTraceIDKey, requestTraceID)
			r = r.WithContext(ctx)

			pathsToIgnoreLogging := []string{
				"/metrics", // to avoid log pollution with /metrics endpoint being hit
			}

			// skip unwanted paths
			if !slices.Contains(pathsToIgnoreLogging, r.URL.RequestURI()) {
				// log the endpoint call metrics
				defer func() {
					log.Info().Fields(map[string]interface{}{
						"method":      r.Method,
						"requester":   r.RemoteAddr,
						"trace-id":    requestTraceID,
						"origin":      r.Header.Get("origin"),
						"path":        r.URL.RequestURI(),
						"req-bytes":   r.ContentLength,
						"resp-status": rw.Status,
						"resp-bytes":  rw.BytesWritten,
						"duration-ms": time.Since(start).Milliseconds(),
					}).Msg(
						"endpoint call",
					)
				}()
			}

			next.ServeHTTP(rw, r)
		})
	}
}
