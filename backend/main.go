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

	"radare-datarecon/backend/internal/budget"
	"radare-datarecon/backend/internal/config"
	"radare-datarecon/backend/internal/database"
	"radare-datarecon/backend/internal/handlers"
	"radare-datarecon/backend/internal/middleware"
	"radare-datarecon/backend/internal/models"
)

func main() {
	// Load application configuration from environment variables.
	cfg := config.Load()

	// Connect to the database and migrate the schema.
	database.Connect(cfg)
	if err := database.DB.AutoMigrate(&models.User{}, &models.Transaction{}, &models.Category{}); err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	// Instantiate the authentication middleware with the JWT secret.
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	// Instantiate the budget service and handlers.
	budgetService := budget.NewService(database.DB)
	budgetHandlers := budget.NewHandlers(budgetService)

	// Register handlers for the API endpoints.
	// Each handler is wrapped with middleware for logging, error handling, and authentication.
	http.Handle("/api/register", middleware.LoggingMiddleware(middleware.ErrorHandler(handlers.Register)))
	http.Handle("/api/login", middleware.LoggingMiddleware(middleware.ErrorHandler(handlers.LoginHandler(cfg.JWTSecret))))
	http.Handle("/healthz", middleware.LoggingMiddleware(middleware.ErrorHandler(handlers.HealthCheck)))

	// Budgeting endpoints
	http.Handle("/api/categories", middleware.LoggingMiddleware(authMiddleware(middleware.ErrorHandler(budgetHandlers.HandleCategories))))
	http.Handle("/api/transactions", middleware.LoggingMiddleware(authMiddleware(middleware.ErrorHandler(budgetHandlers.HandleTransactions))))
	http.Handle("/api/budget/summary", middleware.LoggingMiddleware(authMiddleware(middleware.ErrorHandler(budgetHandlers.GetMonthlySummary))))

	// Create and configure the HTTP server.
	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      http.DefaultServeMux,
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