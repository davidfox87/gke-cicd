package middleware

import (
	"net/http"
	"strings"
)

// Healthcheck returns a middleware that responds to healthcheck requests.
func Healthcheck(handler http.Handler, endpoint string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.Method == "GET" || r.Method == "HEAD") && strings.EqualFold(r.URL.Path, endpoint) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			//nolint errcheck
			w.Write([]byte("Healthy"))
			return
		}

		handler.ServeHTTP(w, r)
	})
}
