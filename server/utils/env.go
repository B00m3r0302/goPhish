package utils

import (
	"log"
	"os"
)

// GetEnv fetches the value of an environment variable or returns a default value if not set
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// MustGetEnv fetches the value of an environment variable or logs a fatal error if not set
func MustGetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Fatalf("Environment variable %s is required but not set", key)
	return ""
}
