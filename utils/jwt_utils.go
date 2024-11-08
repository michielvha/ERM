package utils

import (
    "time"
    "os"
    "github.com/golang-jwt/jwt/v5"
)

// var jwtSecret = []byte("yoursecuresecret") // Replace with a secure secret in production
// Production safe way of handling secret via env var.
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))


// GenerateJWT generates a new JWT token
func GenerateJWT(username string) (string, error) {
    claims := jwt.MapClaims{
        "username": username,
        "exp":      time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}
