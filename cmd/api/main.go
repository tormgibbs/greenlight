package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/emmasela/greenlight/internal/data"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const version = "1.0.0"

// Config struct to hold configuration settings
type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

// Application struct to hold dependencies
type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Initialize a config variable
	var cfg config

	// Set the default values for the config
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(
		&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgreSQL DSN",
	)

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Printf("database connection pool established")

	// Create a new application pointer and assign the config and logger
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	// Declare an HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start HTTP server
	logger.Printf("Starting %s server on port %s", cfg.env, server.Addr)
	err = server.ListenAndServe()
	logger.Fatal(err)
}

// openDB initializes a connection pool to the database using the provided configuration.
func openDB(cfg config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a new connection pool using the context and the database DSN from the configuration.
	db, err := pgxpool.New(ctx, cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify that the connection is successful.
	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

//142
