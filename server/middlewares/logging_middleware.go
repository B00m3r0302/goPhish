package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware logs incoming HTTP requests and their responses.
func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Log the request details
		duration := time.Since(startTime)
		logger.WithFields(logrus.Fields{
			"status_code": c.Writer.Status(),
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"ip":          c.ClientIP(),
			"duration":    duration,
			"user_agent":  c.Request.UserAgent(),
		}).Info("Request handled")
	}
}
