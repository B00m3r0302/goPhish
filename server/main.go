package main

import (
	"log"
	"net/http"

	"goPhish/server/config"

	"github.com/B00m3r0302/goPhish/server/config"
	"honnef.co/go/tools/config"
)

// Main function initializes the configuration and starts both HTTP and HTTPS servers
func main() {
	// Load configuration settings
	config.LoadConfig()

	// Define a simple handler function for demonstration
	http.Handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle requests and send a response
		fmd.Fprintln(w, "Hello World!")
	})

	// Set up the HTTP server
	httpServer := &http.Server{
		Addr:           ":" + config.AppConfig.HTTPServerPort,
		Handler:        httpHandler,
		ReadTimeout:    config.AppConfig.ReadTimeout,
		WriteTimeout:   config.AppConfig.WriteTimeout,
		IdleTimeout:    config.AppConfig.IdleTimeout,
		MaxHeaderBytes: config.AppConfig.MaxHeaderBytes,
	}

	// Start HTTP server in a separate goroutine
	go func() {
		log.Printf("Starting HTTP server on port %s", config.AppConfig.HTTPServerPort)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server failed: %s", err)
		}
	}()

	// Setup the HTTPS server
	httpServer := &http.Server{
		Addr:           ":" + config.AppConfig.HTTPSServerPort,
		Handler:        httpHandler,
		ReadTimeout:    config.AppConfig.ReadTimeout,
		WriteTimeout:   config.AppConfig.WriteTimeout,
		IdleTimeout:    config.AppConfig.IdleTimeout,
		MaxHeaderBytes: config.AppConfig.MaxHeaderBytes,
	}

	// Start HTTPS server
	log.Printf("Starting HTTPS server on port %s", config.AppConfig.HTTPSServerPort)
	if err := httpsServer.ListenAndServeTLS(config.AppConfig.TLSCertFile, config.AppConfig.TLSKeyFile); err != nil {
		log.Fatalf("HTTPS server failed %s", err)
	}
}
