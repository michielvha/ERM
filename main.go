package main

import (
    // only needed to autogen env var, in theory we could keep this behaviour, but allow it to be overwritten by user. in this case if it's not set the binary will still run.
    "crypto/rand"
    "encoding/base64"
//     "log"


    // actually used
    "os"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "erm/db"
    "github.com/gin-gonic/gin"
    "github.com/MKTHEPLUGG/ERM/middleware"
    "github.com/MKTHEPLUGG/ERM/handlers"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    "database/sql"
)

// generateRandomSecret creates a 256-bit (32-byte) random secret
func generateRandomSecret() string {
    secret := make([]byte, 32)
    _, err := rand.Read(secret)
    if err != nil {
        log.Fatal().Msgf("Failed to generate secret: %v", err)
    }
    return base64.StdEncoding.EncodeToString(secret)
}

// setEnvVar sets the environment variable for JWT_SECRET
func setEnvVar() {
    secret := generateRandomSecret()
    err := os.Setenv("JWT_SECRET", secret)
    if err != nil {
        log.Fatal().Msgf("Failed to set environment variable: %v", err)
    }
    log.Info().Msgf("Environment variable JWT_SECRET set: %s", secret) // For testing purposes only; remove this log in production
}

func runMigrations(db *sql.DB) {
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatal("Failed to create migration driver:", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations", // Replace with your migrations directory
        "postgres", driver)
    if err != nil {
        log.Fatal("Failed to create migration instance:", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal("Migration failed:", err)
    } else {
        log.Println("Migrations applied successfully")
    }
}

// function to init our env
func init(){
    // initialize zerologger
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

    // Examples
//     log.Info().Msg("Application started")
//     log.Warn().Msg("This is a warning from someFunction")
//     log.Error().Msg("This is an error from someFunction")

}

func main() {
    // Initialize the database connection
    dbConn, err := db.InitDB()
    if err != nil {
        log.Fatal().Msg("Could not establish a database connection")
    } else {
        log.Info().Msg("Establish database connection")
    }
    defer dbConn.Close() // Ensure that the connection is closed when the program exits

    // Run database migrations
    runMigrations(dbConn)


    // Create a new Gin Router
    r := gin.Default()

    // Public route for login (no JWT middleware)
    r.POST("/login", handlers.Login)

    // Group routes that require authentication
    protected := r.Group("/protected")
    protected.Use(middleware.JWTAuthMiddleware())
    {
        protected.GET("/", handlers.ProtectedEndpoint)
    }

    r.Run("0.0.0.0:8080") // Run the server on port 8080
    log.Info().Msg("Server started on 0.0.0.0:8080")
}