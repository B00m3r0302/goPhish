package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"your_project/server"
	"your_project/server/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize the logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(logrus.InfoLevel)

	// Load configuration using Viper and set up file watcher
	cfg, err := config.LoadConfig("server/config/config.yaml")
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}
	logger.WithField("config", cfg).Info("Configuration loaded successfully")

	// Set up the router using Gin
	router := server.SetupRouter()

	// Determine if the server should run in HTTP or HTTPS mode based on config
	go func() {
		if cfg.Server.UseHTTPS {
			logger.WithField("port", cfg.Server.HTTPSPort).Info("Starting HTTPS server")
			if err := router.RunTLS(":"+cfg.Server.HTTPSPort, cfg.Server.TLSCertFile, cfg.Server.TLSKeyFile); err != nil {
				logger.WithError(err).Fatal("Failed to start HTTPS server")
			}
		} else {
			logger.WithField("port", cfg.Server.HTTPPort).Info("Starting HTTP server")
			if err := router.Run(":" + cfg.Server.HTTPPort); err != nil {
				logger.WithError(err).Fatal("Failed to start HTTP server")
			}
		}
	}()

	// Graceful shutdown handling
	gracefulShutdown(router, logger)
}

// gracefulShutdown handles server shutdown on system interrupt or termination signals
func gracefulShutdown(router *gin.Engine, logger *logrus.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Allow a grace period for ongoing requests to complete
	timeout := time.Second * 10
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	shutdownChan := make(chan struct{})

	go func() {
		if err := router.Shutdown(nil); err != nil {
			logger.WithError(err).Error("Server forced to shutdown")
		}
		close(shutdownChan)
	}()

	select {
	case <-timer.C:
		logger.Warn("Timeout reached, forcing shutdown")
	case <-shutdownChan:
		logger.Info("Server gracefully stopped")
	}
}
