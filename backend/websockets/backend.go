// backend.go - Entry point for the goPhish C2 framework backend
// Author: Hudson Woomer
// Description: This file initializes the backend server, sets up the HTTP server,
// handles routing, middleware, database connections, WebSocket management, and more.

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"goPhish/backend/middlewares"
	"goPhish/backend/services"
	"goPhish/backend/websockets"
	"goPhish/config"

	"github.com/gorilla/mux"
)

// Server configuration
var (
	serverPort      = ":7777" // Default server port
	readTimeout     = 15 * time.Second
	writeTimeout    = 15 * time.Second
	idleTimeout     = 60 * time.Second
	shutdownTimeout = 10 * time.Second
)

// Global variables
var (
	router     *mux.Router        // Router for handling HTTP routes
	wsHub      *websockets.Hub    // WebSocket hub for managing connections
	db         *services.Database // Database service
	shutdownWg sync.WaitGroup     // WaitGroup for graceful shutdown
)

func init() {
	// Load configuration from config.yaml
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the database connection
	db = services.NewDatabase(config.GetDatabaseConfig())
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize WebSocket hub
	wsHub = websockets.NewHub()

	// Initialize the router
	router = mux.NewRouter()

	// Add middlewares
	router.Use(middlewares.LoggingMiddleware)
	router.Use(middlewares.AuthMiddleware)
	router.Use(middlewares.CORSMiddleware)

	// Register API routes
	v1.RegisterRoutes(router, wsHub, db)

	// Serve static files (e.g., for frontend)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
}

// startWebSocketHub starts the WebSocket hub and listens for connections
func startWebSocketHub() {
	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		wsHub.Run()
	}()
}

// startHTTPServer starts the HTTP server with the configured settings
func startHTTPServer() *http.Server {
	server := &http.Server{
		Addr:         serverPort,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	shutdownWg.Add(1)
	go func() {
		defer shutdownWg.Done()
		log.Printf("Starting HTTP server on %s", serverPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	return server
}

// gracefulShutdown handles graceful shutdown of the server, WebSocket hub, and other services
func gracefulShutdown(server *http.Server) {
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	<-shutdownCh

	log.Println("Shutdown signal received. Initiating graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown HTTP server
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Stop WebSocket hub
	wsHub.Shutdown()

	// Close database connection
	if err := db.Close(); err != nil {
		log.Printf("Database close error: %v", err)
	}

	shutdownWg.Wait()
	log.Println("Graceful shutdown complete.")
}

func main() {
	// Start WebSocket hub
	startWebSocketHub()

	// Start HTTP server
	server := startHTTPServer()

	// Handle graceful shutdown
	gracefulShutdown(server)

	log.Println("goPhish C2 framework backend stopped.")
}
