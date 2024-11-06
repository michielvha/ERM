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
    "github.com/gin-gonic/gin"
    "github.com/yourusername/oauth2-jwt-api/middleware"
    "github.com/yourusername/oauth2-jwt-api/handlers"
)

func main() {
    r := gin.Default()

    // Apply JWT authentication middleware to protected routes
    r.Use(middleware.JWTAuthMiddleware())

    // Create a simple protected route
    r.GET("/protected", handlers.ProtectedEndpoint)

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

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secure-secret") // Replace with a secure secret in production

// JWTAuthMiddleware validates the JWT token
func JWTAuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")

        if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
            return
        }

        tokenString = strings.TrimPrefix(tokenString, "Bearer ")

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
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
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secure-secret") // Replace with a secure secret in production

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

### Step 4: Create a Simple Protected Endpoint
**handlers/api_handler.go**:
```go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// ProtectedEndpoint responds only if the request has a valid JWT
func ProtectedEndpoint(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "You have accessed a protected route!"})
}
```

### Step 5: Add an Endpoint to Issue JWT Tokens (Simulating OAuth2)
Create an endpoint to issue tokens to authenticated users for testing purposes.

**main.go (Add)**:
```go
r.POST("/login", handlers.Login)
```

**handlers/api_handler.go (Add)**:
```go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/yourusername/oauth2-jwt-api/utils"
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
```

### Step 6: Test the API
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

### Next Steps
- **Integrate with a real OAuth2 provider** like Auth0 or Keycloak for better security and compliance.
- **Enhance JWT claims** to include user roles, scopes, or additional data.
- **Implement token refresh logic** if necessary.

Would you like any further details on integrating with an OAuth2 provider or refining the JWT authentication flow?