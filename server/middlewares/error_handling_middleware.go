package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorHandlingMiddleware captures any errors that occur during request processing and logs them.
func ErrorHandlingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if errors occurred during request processing
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				// Log the error
				logger.WithError(e.Err).Error("Request error")

				// Send the error response
				c.JSON(http.StatusInternalServerError, gin.H{"error": e.Err.Error()})
			}
			// Abort further processing
			c.Abort()
		}
	}
}
