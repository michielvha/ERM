package handlers

import (
    "database/sql"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/MKTHEPLUGG/ERM/utils"
    "github.com/MKTHEPLUGG/ERM/db"
    "golang.org/x/crypto/bcrypt"
    "github.com/rs/zerolog/log"
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

    dbConn, err := db.InitDB()
    if err != nil {
        log.Error().Msg("Database connection failed")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
        return
    }
    defer dbConn.Close()

    var dbUser struct {
        Username string
        Password string
        Role     string
    }

    // Query the database for the user
    err = dbConn.QueryRow("SELECT username, password, role FROM users WHERE username = $1", user.Username).Scan(&dbUser.Username, &dbUser.Password, &dbUser.Role)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
        }
        return
    }

    // Compare the hashed password with the provided password
    if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    // Generate JWT token
    token, err := utils.GenerateJWT(dbUser.Username, dbUser.Role)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

// ProtectedEndpoint responds only if the request has a valid JWT
func ProtectedEndpoint(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "You have accessed a protected route!"})
}