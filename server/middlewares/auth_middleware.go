package middlewares

import (
	"net/http"
	"strings"

	"your_project/server/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for the presence and validity of the authorization token in the request header.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Split the header to extract the token
		token := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing or malformed"})
			c.Abort()
			return
		}

		// Validate the token
		if valid, err := utils.ValidateToken(token); !valid || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Token is valid, proceed to the next handler
		c.Next()
	}
}
