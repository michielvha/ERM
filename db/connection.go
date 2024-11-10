package db

import (
    "os"
    "fmt"
    "database/sql"
    _ "github.com/lib/pq" // PostgreSQL driver
    "github.com/rs/zerolog/log"
)

func InitDB() (*sql.DB, error) {
    // set connection string to env vars, if none are given use app defaults.
    dbUsername := os.Getenv("POSTGRES_USER")
    if dbUsername == "" {
        dbUsername = "erm_user"
    }

    dbPassword := os.Getenv("POSTGRES_PASSWORD")
    if dbPassword == "" {
        dbPassword = "secure_password"
    }

    dbName := os.Getenv("POSTGRES_DB")
    if dbName == "" {
        dbName = "erm_database"
    }

    // Create connection string
    connStr := fmt.Sprintf("postgresql://%s:%s@db:5432/%s?sslmode=disable", dbUsername, dbPassword, dbName)
    // Connect to the database
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal().Msgf("Failed to connect to the database:", err)
        return nil, err
    }

    // Ping to ensure connection is active
    if err := db.Ping(); err != nil {
        log.Fatal().Msgf("Failed to ping the database:", err)
        return nil, err
    }

    log.Info().Msg("Database connection established")
    return db, nil
}

