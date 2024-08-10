package config

import (
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config struct holds the server's configuration settings
type Config struct {
	HTTPServerPort    string        // Port for the HTTP server
	HTTPSServerPort   string        // Port for the HTTPS server
	ReadTimeout       time.Duration // Maximum duration for reading the entire request, including the body
	WriteTimeout      time.Duration // Maximum duration before timeing out writes of the response
	IdleTimeout       time.Duration // Maximum amount of time to wait for the nex request
	MaxHeaderBytes    int           // Maximum size of request headers in bytes
	JWTSecret         string        // Secret key used for JWT authentication
	DatabaseURL       string        // URL for the database connection
	LogLevel          string        // Level of logging: DEBUG, INFO, WARN, ERROR
	EncryptionKey     string        // Key used for encrypting sensitive data
	AllowedOrigins    []string      // CORS: Allowed origins for the API
	APIVersion        string        // Version of the API
	DefaultAgentGroup string        // Default group for newly registered agents
	TLSCertFile       string        // Path to the TLS certificate file for HTTPS
	TLSKeyFile        string        // Path to the TLS key file for HTTPS
}

// Global variable to hold the loaded configuration
var AppConfig Config

// LoadConfig initializes the configuration by reading environment variables, config file, and setting defaults
func LoadConfig() {
	// Setup default values
	viper.SetDefault("HTTP_SERVER_PORT", "7777")
	viper.SetDefault("HTTPS_SERVER_PORT", "8888")
	viper.SetDefault("READ_TIMEOUT", 15*time.Second)
	viper.SetDefault("WRITE_TIMEOUT", 15*time.Second)
	viper.SetDefault("IDLE_TIMEOUT", 60*time.Second)
	viper.SetDefault("MAX_HEADER_BYTES", 1<<20) // 1 MB
	viper.SetDefault("LOG_LEVEL", "INFO")
	viper.SetDefault("API_VERSION", "v1")
	viper.SetDefault("DEFAULT_AGENT_GROUP", "default")

	// Configure viper to read from environment variables
	viper.AutomaticEnv()

	// Set up configuration file support
	viper.SetConfigName("config") // Name of config file (without extension)
	viper.SetConfigType("yaml")   // Required since the file extension isn't in the name
	viper.AddConfigPath(".")      // Look for config file in the current directory
	viper.ConfigPath("$HOME/.c2") // Optionally look for the config file in the user's home directory
	viper.ConfigPath("/etc/c2")   // Optionally look for the config file in the /etc/c2 directory

	// Read in the config file if available
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Warning: Config file not found: %s", err)
	} else {
		log.Printf("Config file loaded: %s", viper.ConfigFileUsed())
	}

	// Populate the AppConfig struct with the config values
	reloadConfig()

	// Watch for changes in the config file and reload automatically
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed:%s", e.Name)
		reloadConfig()
	})
}

// reloadConfig populates the AppConfig struct with the latest config values
func reloadConfig() {
	AppConfig = Config{
		HTTPServerPort:    viper.GetString("HTTP_SERVER_PORT"),
		HTTPSServerPort:   viper.GetString("HTTPS_SERVER_PORT"),
		ReadTimeout:       viper.GetDuration("READ_TIMEOUT"),
		WriteTimeout:      viper.GetDuration("WRITE_TIMEOUT"),
		IdleTimeout:       viper.GetDuration("IDLE_TIMEOUT"),
		MaxHeaderBytes:    viper.GetInt("MAX_HEADER_BYTES"),
		JWTSecret:         getEnvOrPanic("JWT_SECRET"), // Mandatory setting, panic if not found
		DatabaseURL:       getEnvOrDefault("DATABASE_URL", "postgres://user:password@localhost:5432/c2db?sslmode=disable"),
		LogLevel:          viper.GetString("LOG_LEVEL"),
		EncryptionKey:     getEnvOrPanic("ENCRYPTION_KEY"), // Mandatory setting, panic if not found
		AllowedOrigins:    viper.GetStringSlice("ALLOWED_ORIGINS"),
		APIVersion:        viper.GetString("API_VERSION"),
		DefaultAgentGroup: viper.GetString("DEFAULT_AGENT_GROUP"),
		TLSCertFile:       viper.GetString("TLS_CERT_FILE"),
		TLSKeyFile:        viper.GetString("TLS_KEY_FILE"),
	}

	log.Printf("Config reloaded: %+v", AppConfig)
}

// getEnvOrPanic retrieves an environment variable or panics if it's not set
func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}

// getEnvOrDefault retrieves an environment variable or returns a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
