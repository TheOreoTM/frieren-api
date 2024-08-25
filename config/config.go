package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Configuration holds the application configuration.
type Configuration struct {
	Port        string
	DatabaseURL string
	LogLevel    string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() *Configuration {
	return &Configuration{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "localhost:5432"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv retrieves the value of an environment variable or returns a default value if not set.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// SetupLogger initializes and returns a Logrus logger.
func SetupLogger(level string) *logrus.Logger {
	logger := logrus.New()
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logger.SetLevel(logrus.InfoLevel) // Default to InfoLevel if parsing fails
	} else {
		logger.SetLevel(logLevel)
	}
	return logger
}
