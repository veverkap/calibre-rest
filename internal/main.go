package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const Version = "0.1.0"

func main() {
	// Define command line flags
	var (
		dev        = flag.Bool("dev", false, "Start in dev/debug mode")
		calibreDB  = flag.String("calibre", "", "Path to calibre binary directory")
		library    = flag.String("library", "", "Path to calibre library")
		username   = flag.String("username", "", "Calibre library username")
		password   = flag.String("password", "", "Calibre library password")
		logLevel   = flag.String("log-level", "", "Log level (DEBUG, INFO, WARNING, ERROR)")
		bindAddr   = flag.String("bind", "", "Bind address HOST:PORT")
		version    = flag.Bool("version", false, "Print version")
		help       = flag.Bool("help", false, "Show help")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *version {
		fmt.Printf("v%s\n", Version)
		return
	}

	// Create configuration
	var config *Config
	if *dev {
		config = NewDevConfig()
	} else {
		config = NewProdConfig()
	}

	// Override config with command line flags
	if *calibreDB != "" {
		config.SetCalibreDB(*calibreDB)
	}
	if *library != "" {
		config.SetLibrary(*library)
	}
	if *bindAddr != "" {
		config.SetBindAddr(*bindAddr)
	}
	if *username != "" {
		config.SetUsername(*username)
	}
	if *password != "" {
		config.SetPassword(*password)
	}
	if *logLevel != "" {
		config.SetLogLevel(*logLevel)
	}

	// Setup logger
	logger := log.New(os.Stdout, "[calibre-rest] ", log.LstdFlags)
	if config.Debug {
		logger.Printf("Server config: %+v", config)
	}

	// Create Calibre wrapper
	calibre := NewCalibreWrapper(config.CalibreDB, config.Library, config.Username, config.Password, logger)
	
	// Check calibre setup
	if err := calibre.Check(); err != nil {
		logger.Fatalf("Failed to initialize calibre wrapper: %v", err)
	}

	// Create server
	server := NewServer(calibre, logger, Version)
	router := server.SetupRoutes()

	// Setup HTTP server
	httpServer := &http.Server{
		Addr:    config.BindAddr,
		Handler: router,
		// Good practice: enforce timeouts
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Printf("Starting server on %s", config.BindAddr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exited")
}