package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

// Config struct to hold configuration settings
type config struct {
	port int
	env string
}

// Application struct to hold dependencies
type application struct {
	config config
	logger *log.Logger
}


func main() {
	// Initialize a config variable
	var cfg config

	// Set the default values for the config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Create a new application pointer and assign the config and logger
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Declare a new servermux and register the routes and handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthCheckHandler)

	// Declare an HTTP server
	server := &http.Server{
		Addr: fmt.Sprintf("%d", cfg.port),
		Handler: mux,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start HTTP server
	logger.Printf("Starting server on port %d in %s mode", cfg.port, cfg.env)
	err := server.ListenAndServe()
	logger.Fatal(err)
}