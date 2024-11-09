package db

import (
    "database/sql"
    _ "github.com/lib/pq" // PostgreSQL driver
    "github.com/rs/zerolog/log"
)

func InitDB() (*sql.DB, error) {
    connStr := "postgresql://admin:secure_admin_password@db:5432/erm_database?sslmode=disable" // TODO: remove hardcoding from var to env var.
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

