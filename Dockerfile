# Stage 1: Build the Go application
FROM golang:1.22-alpine3.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies are cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app with optimizations
RUN go build -o main .

# Stage 2: Create a smaller image for running the Go application
FROM alpine:3.20.3

# Install certificates for HTTPS if needed
RUN apk --no-cache add ca-certificates

# Create a non-root user with a specific UID
RUN addgroup -g 1000 appgroup && adduser -u 1000 -S appuser -G appgroup

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the builder
COPY --from=builder /app/main .

# Change ownership of the binary to the non-root user
RUN chown appuser:appgroup /app/main

# Switch to the non-root user
USER appuser

# Command to run the executable
CMD ["./main"]
