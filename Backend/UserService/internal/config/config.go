package config

import (
    "fmt"
    "log"
    "os"
)

// Config holds all DB connection info
type Config struct {
    DatabaseURL string
    JWTSecret   string
    Port        string
}

func Load() Config {
    user := os.Getenv("POSTGRES_USER")
    pass := os.Getenv("POSTGRES_PASS")
    host := os.Getenv("POSTGRES_HOST")
    port := os.Getenv("POSTGRES_PORT")
    db   := os.Getenv("POSTGRES_DB")
    ssl  := os.Getenv("SSL_MODE")
    jwt := os.Getenv("JWT_SECRET")
    if user == "" || pass == "" || host == "" || port == "" || db == "" {
        log.Fatal("Database credentials must be set")
    }

	if ssl == "" {
		ssl = "disable"
	}

	if host == "" {
		host = "localhost"
	}

    // Build the Database URL
    databaseURL := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        user, pass, host, port, db, ssl,
    )

    httpPort := os.Getenv("PORT")
    if httpPort == "" {
        httpPort = "8080"
    }

    return Config{
        DatabaseURL: databaseURL,
        JWTSecret:   jwt,
        Port:        httpPort,
    }
}
