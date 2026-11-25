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
	"car-maintenance/internal/database"
	"car-maintenance/internal/handler"
	"car-maintenance/internal/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Invalid configuration", "error", err)
		os.Exit(1)
	}

	// Initialize SQLite database
	db, err := database.Connect(cfg.DBPath)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err, "path", cfg.DBPath)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Error closing database connection", "error", err)
		} else {
			slog.Info("Database connection closed cleanly")
		}
	}()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		slog.Error("Failed to run database migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("Database migrations completed successfully", "path", cfg.DBPath)

	// Handlers and middlewares
	isProd := cfg.Env == "production"
	authHandler := handler.NewAuthHandler(cfg.AppAccessKey, cfg.JWTSecret, isProd)
	vehicleHandler := handler.NewVehicleHandler(db)
	maintenanceTypeHandler := handler.NewMaintenanceTypeHandler(db)
	vehicleMaintenanceHandler := handler.NewVehicleMaintenanceHandler(db)
	dashboardHandler := handler.NewDashboardHandler(db)
	requireAuth := middleware.RequireAuth(cfg.JWTSecret)
	loginRateLimiter := middleware.NewRateLimiter(5, time.Minute)

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /api/health", handler.HealthHandler)
	mux.Handle("POST /api/login", loginRateLimiter.Middleware(http.HandlerFunc(authHandler.Login)))
	mux.HandleFunc("POST /api/logout", authHandler.Logout)

	// Protected routes - Auth
	mux.Handle("GET /api/auth/verify", requireAuth(http.HandlerFunc(authHandler.Verify)))

	// Protected routes - Vehicles
	mux.Handle("GET /api/vehicles", requireAuth(http.HandlerFunc(vehicleHandler.List)))
	mux.Handle("POST /api/vehicles", requireAuth(http.HandlerFunc(vehicleHandler.Create)))
	mux.Handle("PUT /api/vehicles/{id}", requireAuth(http.HandlerFunc(vehicleHandler.Update)))
	mux.Handle("PATCH /api/vehicles/{id}/mileage", requireAuth(http.HandlerFunc(vehicleHandler.UpdateMileage)))
	mux.Handle("DELETE /api/vehicles/{id}", requireAuth(http.HandlerFunc(vehicleHandler.Delete)))

	// Protected routes - Maintenance Types
	mux.Handle("GET /api/maintenance-types", requireAuth(http.HandlerFunc(maintenanceTypeHandler.List)))
	mux.Handle("POST /api/maintenance-types", requireAuth(http.HandlerFunc(maintenanceTypeHandler.Create)))
	mux.Handle("PUT /api/maintenance-types/{id}", requireAuth(http.HandlerFunc(maintenanceTypeHandler.Update)))
	mux.Handle("DELETE /api/maintenance-types/{id}", requireAuth(http.HandlerFunc(maintenanceTypeHandler.Delete)))

	// Protected routes - Vehicle Maintenances
	mux.Handle("GET /api/vehicles/{id}/maintenances", requireAuth(http.HandlerFunc(vehicleMaintenanceHandler.ListByVehicle)))
	mux.Handle("POST /api/vehicles/{id}/maintenances", requireAuth(http.HandlerFunc(vehicleMaintenanceHandler.Create)))
	mux.Handle("PUT /api/vehicle-maintenances/{id}", requireAuth(http.HandlerFunc(vehicleMaintenanceHandler.Update)))
	mux.Handle("DELETE /api/vehicle-maintenances/{id}", requireAuth(http.HandlerFunc(vehicleMaintenanceHandler.Delete)))

	// Protected routes - Dashboard
	mux.Handle("GET /api/dashboard", requireAuth(http.HandlerFunc(dashboardHandler.Get)))

	// Chain global middlewares: Logger, SecurityHeaders, and CORS
	corsMiddleware := middleware.CORS(cfg.CORSAllowedOrigins)
	securityMiddleware := middleware.SecurityHeaders(isProd)
	handlerWithMiddlewares := middleware.Logger(securityMiddleware(corsMiddleware(mux)))

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
