# Use Go 1.25 to match go.mod requirement
FROM golang:1.25-alpine

# Set working directory inside the container
WORKDIR /app

# Install necessary packages for building Go projects
RUN apk add --no-cache git bash

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go application
RUN go build -o main ./cmd/main.go

# Expose port 8080 for the API
EXPOSE 8080

# Command to run the application
CMD ["./main"]
