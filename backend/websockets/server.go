// server.go - Server management for the goPhish C2 framework
// Author: Hudson Woomer
// Description: This file contains functions and configurations for managing
// the HTTP server, including setup, start, and shutdown procedures.

package backend

import (
	"context"
	"log"
	"net/http"
	"sync"
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

// Server holds the HTTP server instance and its dependencies
type Server struct {
	httpServer *http.Server
	router     *mux.Router
	wsHub      *websockets.Hub
	db         *services.Database
	shutdownWg sync.WaitGroup
}

// NewServer initializes a new server with the necessary configurations
func NewServer() *Server {
	// Load configuration from config.yaml
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the database connection
	db := services.NewDatabase(config.GetDatabaseConfig())
	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize WebSocket hub
	wsHub := websockets.NewHub()

	// Initialize the router
	router := mux.NewRouter()

	// Add middlewares
	router.Use(middlewares.LoggingMiddleware)
	router.Use(middlewares.AuthMiddleware)
	router.Use(middlewares.CORSMiddleware)

	// Register API routes
	v1.RegisterRoutes(router, wsHub, db)

	// Serve static files (e.g., for frontend)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Create and return the Server instance
	return &Server{
		httpServer: &http.Server{
			Addr:         serverPort,
			Handler:      router,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
		router: router,
		wsHub:  wsHub,
		db:     db,
	}
}

// Start begins the HTTP server and WebSocket hub
func (s *Server) Start() {
	s.shutdownWg.Add(1)
	go func() {
		defer s.shutdownWg.Done()
		log.Printf("Starting HTTP server on %s", serverPort)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	s.shutdownWg.Add(1)
	go func() {
		defer s.shutdownWg.Done()
		s.wsHub.Run()
	}()
}

// Shutdown handles graceful shutdown of the server, WebSocket hub, and other services
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Println("Shutdown signal received. Initiating graceful shutdown...")

	// Shutdown HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Stop WebSocket hub
	s.wsHub.Shutdown()

	// Close database connection
	if err := s.db.Close(); err != nil {
		log.Printf("Database close error: %v", err)
	}

	s.shutdownWg.Wait()
	log.Println("Graceful shutdown complete.")
}
