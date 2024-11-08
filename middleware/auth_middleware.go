package middleware

import (
    "net/http"
    "strings"
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

// Replace with a secure secret in production
// var jwtSecret = []byte("yoursecuresecret")

// Production safe way of handling secret via env var.
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))


// JWTAuthMiddleware: validates the JWT token
func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get the Authorization header
        tokenString := c.GetHeader("Authorization")

        // Check if the token is present and starts with "Bearer "
        if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
//            log.Printf("Missing or malformed token error: %v", tokenString)  // only for debug, exposes token serverside, bad practise.
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
            return
        }

        // Extract the actual token part, removing the "Bearer " prefix
        tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))

        // Parse the token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Validate the signing method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return jwtSecret, nil
        })

        // Handle token parsing errors
        if err != nil {
            log.Printf("Token parsing error: %v", err) // Logs internally; do not expose this to clients
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or malformed token"})
            return
        }

        if !token.Valid {
            log.Printf("Token validation failed: token is not valid") // Logs when the token is invalid
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        // Optionally, set the token claims to context for further use
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            c.Set("claims", claims)
        } else {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            return
        }

        c.Next()
    }
}
