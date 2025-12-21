// Package handler provides the HTTP handler with all routes and middleware configured.
package handler

import (
	"log/slog"
	"net/http"

	"github.com/BadrChoubai/hello-api/pkg/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewHandler returns an HTTP handler with all routes and middleware configured.
func NewHandler(log *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux)

	var handler http.Handler = mux

	// Compose middleware.
	handler = middleware.Healthcheck(handler, "/healthz")
	handler = middleware.Logging(handler, log)

	// Instrument the entire server.
	return otelhttp.NewHandler(handler, "/")
}

func addRoutes(mux *http.ServeMux) {
	// handleFunc wraps mux.HandleFunc and enriches the handler's
	// HTTP instrumentation with the pattern as the http.route.
	handle := func(pattern string, h http.Handler) {
		h = otelhttp.WithRouteTag(pattern, h)
		mux.Handle(pattern, h)
	}

	handle("/", http.NotFoundHandler())
	handle("/metrics", promhttp.Handler())

	handle("/greet", HandleGreeting())
}
