package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Logging returns a middleware that logs basic HTTP request information.
func Logging(handler http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		handler.ServeHTTP(w, r)

		logger.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.Query(),
			"duration", time.Since(start),
			"remote_addr", r.RemoteAddr,
		)
	})
}
