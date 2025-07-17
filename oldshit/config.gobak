package main

import (
	"log"
	"os"
	"strings"
)

// Config holds application configuration
type Config struct {
	CalibreDB string
	Library   string
	BindAddr  string
	Username  string
	Password  string
	LogLevel  string
	Debug     bool
	Testing   bool
}

// Valid log levels
var validLogLevels = []string{"DEBUG", "INFO", "WARNING", "ERROR"}

// NewConfig creates a new configuration with defaults from environment variables
func NewConfig() *Config {
	return &Config{
		CalibreDB: getEnv("CALIBRE_REST_PATH", "/opt/calibre/calibredb"),
		Library:   getEnv("CALIBRE_REST_LIBRARY", "/library"),
		BindAddr:  getEnv("CALIBRE_REST_ADDR", "localhost:5000"),
		Username:  getEnv("CALIBRE_REST_USERNAME", ""),
		Password:  getEnv("CALIBRE_REST_PASSWORD", ""),
		LogLevel:  getEnv("CALIBRE_REST_LOG_LEVEL", "INFO"),
		Debug:     false,
		Testing:   false,
	}
}

// NewDevConfig creates a development configuration
func NewDevConfig() *Config {
	config := NewConfig()
	config.LogLevel = getEnv("CALIBRE_REST_LOG_LEVEL", "DEBUG")
	config.Debug = true
	return config
}

// NewTestConfig creates a test configuration
func NewTestConfig() *Config {
	config := NewConfig()
	config.Testing = true
	return config
}

// NewProdConfig creates a production configuration
func NewProdConfig() *Config {
	return NewConfig()
}

// SetCalibreDB sets the calibredb path
func (c *Config) SetCalibreDB(path string) {
	if path != "" {
		c.CalibreDB = path
	}
}

// SetLibrary sets the library path
func (c *Config) SetLibrary(path string) {
	if path != "" {
		c.Library = path
	}
}

// SetBindAddr sets the bind address
func (c *Config) SetBindAddr(addr string) {
	if addr != "" {
		c.BindAddr = addr
	}
}

// SetUsername sets the username
func (c *Config) SetUsername(username string) {
	c.Username = username
}

// SetPassword sets the password
func (c *Config) SetPassword(password string) {
	c.Password = password
}

// SetLogLevel sets the log level with validation
func (c *Config) SetLogLevel(level string) {
	if level == "" {
		return
	}
	
	level = strings.ToUpper(level)
	if !contains(validLogLevels, level) {
		log.Printf("Log level %q not supported. Setting log level to INFO", level)
		c.LogLevel = "INFO"
		return
	}
	
	c.LogLevel = level
}

// GetBindHost returns the host part of the bind address
func (c *Config) GetBindHost() string {
	parts := strings.Split(c.BindAddr, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return "localhost"
}

// GetBindPort returns the port part of the bind address
func (c *Config) GetBindPort() string {
	parts := strings.Split(c.BindAddr, ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return "5000"
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}