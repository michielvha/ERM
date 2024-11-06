package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/MKTHEPLUGG/ERM/utils"
)

// Login simulates user login and returns a JWT token
func Login(c *gin.Context) {
    var user struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Basic username/password check (for demonstration only; use a secure check in production)
    if user.Username == "testuser" && user.Password == "password" {
        token, err := utils.GenerateJWT(user.Username)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"token": token})
    } else {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
    }
}


// ProtectedEndpoint responds only if the request has a valid JWT
func ProtectedEndpoint(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "You have accessed a protected route!"})
}
