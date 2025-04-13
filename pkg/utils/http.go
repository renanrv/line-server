package utils

import "net/http"

type responseWriter struct {
	ResponseWriter http.ResponseWriter
	Status         int
	BytesWritten   int
}

func (rw *responseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.Status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.BytesWritten += n
	return n, err
}

// WrapResponseWriter wraps the http response wrapper if not wrapped
func WrapResponseWriter(w http.ResponseWriter) (sw *responseWriter) {
	var ok bool
	if sw, ok = w.(*responseWriter); !ok {
		sw = &responseWriter{ResponseWriter: w}
	}
	return sw
}
