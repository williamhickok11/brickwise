package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"brickwise/backend/internal/db"
	"brickwise/backend/internal/handlers"
	"brickwise/backend/internal/logger"
	"brickwise/backend/internal/qbo"
	"brickwise/backend/internal/repository"
	"brickwise/backend/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Initialize logger
	log := logger.Default()
	log.Info("Starting Brickwise API server")

	// Initialize database
	database, err := db.InitDB()
	if err != nil {
		log.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	// Initialize layers
	propertyRepo := repository.NewPropertyRepository(database)
	propertyService := service.NewPropertyService(propertyRepo)
	timeEntryRepo := repository.NewTimeEntryRepository(database)
	timeEntryService := service.NewTimeEntryService(timeEntryRepo, propertyRepo)
	residenceRepo := repository.NewResidenceRepository(database)
	residenceService := service.NewResidenceService(residenceRepo, propertyRepo)
	mortgageRepo := repository.NewMortgageRepository(database)
	mortgageService := service.NewMortgageService(mortgageRepo, propertyRepo)

	qboCfg := qbo.LoadConfig()
	qboClient := qbo.NewClient(qboCfg)
	qboRepo := repository.NewQBORepository(database)
	qboService := service.NewQBOService(qboRepo, propertyRepo, qboCfg, qboClient)
	dashboardService := service.NewDashboardService(propertyRepo, mortgageRepo, residenceRepo, qboService)

	if qboCfg.Configured() {
		log.Info("QuickBooks Online integration configured", "environment", qboCfg.Environment)
	} else {
		log.Info("QuickBooks Online not configured (set QBO_CLIENT_ID and QBO_CLIENT_SECRET)")
	}

	// Setup router
	r := chi.NewRouter()

	// Middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Configure appropriately for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		handlers.RegisterRoutes(r, propertyService, timeEntryService, residenceService, mortgageService, dashboardService, qboService)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Info("Server starting", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("Server exited gracefully")
}
