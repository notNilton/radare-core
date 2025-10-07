// main.go is the entry point for the backend server.
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"radare-datarecon/backend/internal/config"
	"radare-datarecon/backend/internal/database"
	"radare-datarecon/backend/internal/handlers"
	"radare-datarecon/backend/internal/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Load application configuration from environment variables.
	cfg := config.Load()

	// Connect to the database and auto-migrate schemas.
	database.Connect(cfg)

	// Set up the router.
	r := mux.NewRouter()

	// Instantiate the authentication middleware with the JWT secret.
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Public routes
	r.Handle("/api/register", middleware.LoggingMiddleware(middleware.ErrorHandler(handlers.Register))).Methods("POST")
	r.Handle("/api/login", middleware.LoggingMiddleware(middleware.ErrorHandler(handlers.LoginHandler(cfg.JWTSecret)))).Methods("POST")
	r.Handle("/healthz", middleware.LoggingMiddleware(middleware.ErrorHandler(handlers.HealthCheck))).Methods("GET")

	// User-related routes (protected)
	userRouter := r.PathPrefix("/api/user").Subrouter()
	userRouter.Use(authMiddleware)
	userRouter.Handle("/profile", middleware.LoggingMiddleware(middleware.ErrorHandler(handlers.GetUserProfile))).Methods("GET")
	userRouter.Handle("/profile", middleware.LoggingMiddleware(middleware.ErrorHandler(handlers.UpdateUserProfile))).Methods("PUT")
	userRouter.Handle("/password", middleware.LoggingMiddleware(middleware.ErrorHandler(handlers.ChangePassword))).Methods("POST")

	// Fueling-related routes (protected)
	fuelingRouter := r.PathPrefix("/api/fueling").Subrouter()
	fuelingRouter.Use(authMiddleware)
	fuelingRouter.Handle("", middleware.LoggingMiddleware(middleware.ErrorHandler(handlers.CreateFuelingRecord))).Methods("POST")
	fuelingRouter.Handle("", middleware.LoggingMiddleware(middleware.ErrorHandler(handlers.GetUserFuelingRecords))).Methods("GET")
	fuelingRouter.Handle("/{id}", middleware.LoggingMiddleware(middleware.ErrorHandler(handlers.UpdateFuelingRecord))).Methods("PUT")
	fuelingRouter.Handle("/{id}", middleware.LoggingMiddleware(middleware.ErrorHandler(handlers.DeleteFuelingRecord))).Methods("DELETE")


	// Create and configure the HTTP server.
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      r, // Use the gorilla/mux router
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Set up a channel to listen for OS signals (SIGINT, SIGTERM) for graceful shutdown.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a separate goroutine.
	go func() {
		log.Printf("Server starting on port %s...", cfg.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	// Block until a shutdown signal is received.
	sig := <-sigChan
	log.Printf("Received shutdown signal: %v, initiating graceful shutdown...\n", sig)

	// Create a context with a timeout to allow for graceful shutdown.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server.
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}

	log.Println("Server shut down successfully.")
}