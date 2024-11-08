package db

import (
    "database/sql"
    _ "github.com/lib/pq" // PostgreSQL driver
    "log"
)

func InitDB() (*sql.DB, error) {
    connStr := "postgresql://erm_user:secure_password@db:5432/erm_database?sslmode=disable" // TODO: remove hardcoding from var to env var.
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
        return nil, err
    }

    // Ping to ensure connection is active
    if err := db.Ping(); err != nil {
        log.Fatal("Failed to ping the database:", err)
        return nil, err
    }

    log.Println("Database connection established")
    return db, nil
}