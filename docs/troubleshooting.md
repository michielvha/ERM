# Troubleshooting Guide

Details on how to troubleshoot the API.

## Linux environment

1. **Run the server**:
   ```bash
   go get 
   go run main.go
   ```

2. **Login to get a JWT token**:
   ```bash
   curl -X POST http://localhost:8080/login -d '{"username": "admin", "password": "secure_admin_password"}' -H "Content-Type: application/json"
   ```


3. **Access the protected route**:
   ```bash
   curl -H "Authorization: Bearer <your-jwt-token>" http://localhost:8080/protected
   ```

## Windows environment

**Powershell Example**:
   ````powershell
   $response = Invoke-RestMethod -Method Post -Uri http://localhost:8080/login -Body '{"username": "admin", "password": "secure_admin_password"}' -ContentType "application/json"
   $token = $response.token
   $token
   Invoke-RestMethod -Method Get -Uri http://localhost:8080/v1/protected -Headers @{ Authorization = "Bearer $token" }
   ````