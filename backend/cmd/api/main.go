package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"car-maintenance/internal/config"
	"car-maintenance/internal/handler"
	"car-maintenance/internal/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()

	mux := http.NewServeMux()

	// Register routes using Go standard library net/http (method + path routing)
	mux.HandleFunc("GET /api/health", handler.HealthHandler)

	// Chain global middlewares: Logger and CORS
	corsMiddleware := middleware.CORS(cfg.CORSAllowedOrigins)
	handlerWithMiddlewares := middleware.Logger(corsMiddleware(mux))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      handlerWithMiddlewares,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server run context for graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, cancel := context.WithTimeout(serverCtx, 15*time.Second)
		defer cancel()

		slog.Info("Shutting down API server gracefully...")
		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				slog.Error("Graceful shutdown timed out... forcing exit.")
				os.Exit(1)
			}
		}()

		if err := server.Shutdown(shutdownCtx); err != nil {
			slog.Error("Server shutdown failed", "error", err)
		}
		serverStopCtx()
	}()

	slog.Info("Backend API server starting", "port", cfg.Port, "env", cfg.Env)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}

	<-serverCtx.Done()
	slog.Info("Backend API server stopped successfully")
}
