## Simple API with JWT & OAuth

Starting with a simplified scope focusing on authentication and basic request handling with OAuth2 and JWT is a smart approach. Here’s how you can set up a basic API in Golang using **Gin** and **OAuth2/JWT** for authentication:

### Step-by-Step Guide to Create an OAuth2 and JWT-Based API in Golang

1. **Set Up Your Project**:
   - Create your project directory:
     ```bash
     mkdir oauth2-jwt-api
     cd oauth2-jwt-api
     ```
   - Initialize the Go module:
     ```bash
     go mod init github.com/yourusername/oauth2-jwt-api
     ```

2. **Install Necessary Dependencies**:
   - Install **Gin**, **jwt-go** (for JWT handling), and **oauth2** packages:
     ```bash
     go get -u github.com/gin-gonic/gin
     go get -u github.com/golang-jwt/jwt/v5
     go get -u golang.org/x/oauth2
     ```

3. **Basic Project Structure**:
   ```
   oauth2-jwt-api/
   ├── main.go
   ├── middleware/
   │   └── auth_middleware.go
   ├── handlers/
   │   └── api_handler.go
   ├── utils/
   │   └── jwt_utils.go
   ```

### Step 1: Create the Server with JWT Middleware
**main.go**:
```go
package main

import (
    // only needed to autogen env var
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

    r.Run(":8080") // Run the server on port 8080
}
```

### Step 2: Create JWT Authentication Middleware
**middleware/auth_middleware.go**:
```go
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

```

### Step 3: Create Utility Functions for JWT Generation
**utils/jwt_utils.go**:
```go
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

```

### Step 4: Test the API
1. **Run the server**:
   ```bash
   go get 
   go run main.go
   ```

2. **Login to get a JWT token**:
   ```bash
   curl -X POST http://localhost:8080/login -d '{"username": "testuser", "password": "password"}' -H "Content-Type: application/json"
   # powershell: curl -Method Post -Uri http://localhost:8080/login -Body '{"username": "testuser", "password": "password"}' -ContentType "application/json"
   ```


3. **Access the protected route**:
   ```bash
   curl -H "Authorization: Bearer <your-jwt-token>" http://localhost:8080/protected
   ```
   
4. **Full Powershell Example**:
   ````powershell
   $response = Invoke-RestMethod -Method Post -Uri http://localhost:8080/login -Body '{"username": "testuser", "password": "password"}' -ContentType "application/json"
   $token = $response.token
   Invoke-RestMethod -Method Get -Uri http://localhost:8080/protected -Headers @{ Authorization = "Bearer $token" }
   ````