package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	healthPath = "/api/v1/health"
)

type statusCodeWrapperResponseWriter struct {
	writer     http.ResponseWriter
	statusCode int
}

func newStatusCodeWrapperResponseWriter(w http.ResponseWriter) *statusCodeWrapperResponseWriter {
	return &statusCodeWrapperResponseWriter{
		writer:     w,
		statusCode: http.StatusOK,
	}
}

// WriteHeader overrides http.ResponseWriter.WriteHeader : records statusCode
func (w *statusCodeWrapperResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.writer.WriteHeader(code)
}

// AccessLog logging when receive request exclude health route
func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newStatusCodeWrapperResponseWriter(w)

		// Serve
		next.ServeHTTP(w, r)

		if r.URL.Path != healthPath {
			elapsed := time.Since(start)
			code := lrw.statusCode
			entry := log.WithFields(log.Fields{
				"path":        r.URL.Path,
				"method":      r.Method,
				"statusCode":  code,
				"elapsedTime": elapsed,
			})
			if code >= http.StatusInternalServerError {
				entry.Error("ACCESS RECEIVED")
			} else if code >= http.StatusBadRequest {
				entry.Warn("ACCESS RECEIVED")
			} else {
				entry.Info("ACCESS RECEIVED")
			}
		}
	})
}
