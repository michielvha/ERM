package main

import (
    "github.com/gin-gonic/gin"
    "github.com/MKTHEPLUGG/ERM/middleware"
    "github.com/MKTHEPLUGG/ERM/handlers"
)

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

    r.Run(":8080") // Run the server on port 8080
}