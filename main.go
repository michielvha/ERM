package main

import (
    // only needed to autogen env var, in theory we could keep this behaviour, but allow it to be overwritten by user. in this case if it's not set the binary will still run.
    "crypto/rand"
    "encoding/base64"
    "log"
    "os"

    // actually used
    "github.com/gin-gonic/gin"
    "github.com/MKTHEPLUGG/ERM/middleware"
    "github.com/MKTHEPLUGG/ERM/handlers"
)

// generateRandomSecret creates a 256-bit (32-byte) random secret
func generateRandomSecret() string {
    secret := make([]byte, 32)
    _, err := rand.Read(secret)
    if err != nil {
        log.Fatalf("Failed to generate secret: %v", err)
    }
    return base64.StdEncoding.EncodeToString(secret)
}

// setEnvVar sets the environment variable for JWT_SECRET
func setEnvVar() {
    secret := generateRandomSecret()
    err := os.Setenv("JWT_SECRET", secret)
    if err != nil {
        log.Fatalf("Failed to set environment variable: %v", err)
    }
    log.Printf("Environment variable JWT_SECRET set: %s", secret) // For testing purposes only; remove this log in production
}

func main() {
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
}