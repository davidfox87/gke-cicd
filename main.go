// Package main implements a simple HTTP server.
package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/BadrChoubai/hello-api/internal/handler"
	"github.com/BadrChoubai/hello-api/internal/observability/logging"
	"github.com/BadrChoubai/hello-api/internal/observability/telemetry"
)

func run(ctx context.Context, stdout io.Writer) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := logging.NewStdLogger(stdout)

	otelShutdown, err := telemetry.SetupOTelSDK(ctx)
	if err != nil {
		logger.Error("initializing OpenTelemetry", "error", err)
		return err
	}
	defer func() {
		err = errors.Join(err, otelShutdown(ctx))
	}()

	server := &http.Server{
		Addr:         ":" + strconv.Itoa(8080),
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      handler.NewHandler(logger),
	}

	srvErr := make(chan error, 1)
	go func() {
		logger.Info("server started", "addr", server.Addr)
		srvErr <- server.ListenAndServe()
	}()

	select {
	case err = <-srvErr:
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", "error", err)
			return err
		}
		logger.Info("server closed")
	case <-ctx.Done():
		logger.Warn("shutdown signal received")
		stop()
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
		logger.Error("server shutdown failed", "error", shutdownErr)
		return shutdownErr
	}

	logger.Info("server shutdown complete")
	return nil
}

func main() {
	ctx := context.Background()

	if err := run(ctx, os.Stdout); err != nil {
		log.Fatalln(err)
	}
}
