package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

// Replace with a secure secret in production
var jwtSecret = []byte("yoursecuresecret")

// JWTAuthMiddleware validates the JWT token
func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get the Authorization header
        tokenString := c.GetHeader("Authorization")

        // Check if the token is present and starts with "Bearer "
        if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
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
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or malformed token"})
            return
        }

        if !token.Valid {
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
