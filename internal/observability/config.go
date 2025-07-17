package observability

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config holds application configuration
type Config struct {
	CalibreDBPath string
	BindAddr      string
	LogLevel      zapcore.Level
	Debug         bool
}

// NewConfig creates a new configuration with defaults from environment variables
func NewConfig() *Config {
	logLevel := zap.InfoLevel // Default log level
	// Read log level from environment variable, defaulting to InfoLevel
	level := getEnv("LOG_LEVEL", "INFO")
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
		CalibreDBPath: getEnv("PATH", "metadata.db"),
		LogLevel:      logLevel,
		Debug:         false,
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
