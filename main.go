package main

import (
    // Only needed to autogen env var
    "crypto/rand"
    "encoding/base64"

    // Used for embedding files
    "embed"
    "os"
    "database/sql"

    // Logging library
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"

    // Project-specific packages
    "github.com/MKTHEPLUGG/ERM/db"
    "github.com/gin-gonic/gin"
    "github.com/MKTHEPLUGG/ERM/middleware"
    "github.com/MKTHEPLUGG/ERM/handlers"

    // Migration libraries
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    "github.com/golang-migrate/migrate/v4/source/iofs"
    _ "github.com/lib/pq"
)
// Embed the `migrations` directory and its content
//go:embed migrations/*.sql
var migrationFiles embed.FS

// function to init our env
func init(){
    // initialize zerologger
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
    log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

//     // Examples
//     log.Info().Msg("Application started")
//     log.Warn().Msg("This is a warning from someFunction")
//     log.Error().Msg("This is an error from someFunction")

}

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
//    log.Info().Msgf("Environment variable JWT_SECRET set: %s", secret) // For testing purposes only; remove this log in production
}

func runMigrations(db *sql.DB) {
    log.Info().Msg("Starting migrations")

    // Create a new iofs driver for the embedded migrations
    sourceDriver, err := iofs.New(migrationFiles, "migrations")
    if err != nil {
        log.Fatal().Msgf("Failed to create iofs source driver: %v", err)
    }

    // Create a database driver
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatal().Msgf("Failed to create migration driver: %v", err)
    }

    // Create the migration instance using the iofs source and database driver
    m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
    if err != nil {
        log.Fatal().Msgf("Failed to create migration instance: %v", err)
    }

    // Run the migrations
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal().Msgf("Migration failed: %v", err)
    } else {
        log.Info().Msg("Migrations applied successfully")
    }
}

func main() {
    // Run database migrations
    runMigrations(dbConn)

    // Initialize the database connection
    dbConn, err := db.InitDB()
    if err != nil {
        log.Fatal().Msg("Could not establish a database connection")
    } else {
        log.Info().Msg("Establish database connection")
    }
    defer dbConn.Close() // Ensure that the connection is closed when the program exits

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