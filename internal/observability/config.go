package observability

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds application configuration
type Config struct {
	CalibreDB string
	Library   string
	BindAddr  string
	Username  string
	Password  string
	LogLevel  zapcore.Level
	Debug     bool
}

// NewConfig creates a new configuration with defaults from environment variables
func NewConfig() *Config {
	logLevel := zap.InfoLevel // Default log level
	// Read log level from environment variable, defaulting to InfoLevel
	level := getEnv("CALIBRE_REST_LOG_LEVEL", "INFO")
	switch strings.ToLower(level) {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warning":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	}

	return &Config{
		CalibreDB: getEnv("CALIBRE_REST_PATH", "/opt/calibre/calibredb"),
		Library:   getEnv("CALIBRE_REST_LIBRARY", "/library"),
		BindAddr:  getEnv("CALIBRE_REST_ADDR", "localhost:5000"),
		Username:  getEnv("CALIBRE_REST_USERNAME", ""),
		Password:  getEnv("CALIBRE_REST_PASSWORD", ""),
		LogLevel:  logLevel,
		Debug:     false,
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
