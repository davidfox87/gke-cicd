// Package logging provides a simple structured logger using slog.
package logging

import (
	"io"
	"log/slog"
)

// NewStdLogger returns a new slog.Logger that writes JSON logs to the provided writer.
func NewStdLogger(w io.Writer) *slog.Logger {
	handler := slog.NewJSONHandler(w, nil)

	return slog.New(handler)
}
