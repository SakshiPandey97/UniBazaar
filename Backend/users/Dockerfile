# Use official Golang image as a builder
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application files
COPY . .

# Build the Go binary for Linux
RUN GOOS=linux GOARCH=amd64 go build -o main .

# Use a minimal Linux base image for deployment
FROM alpine:latest

# Set the working directory in the container
WORKDIR /root/

# Install libc6-compat (needed for some Go binaries)
RUN apk add --no-cache libc6-compat

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/main .

# Explicitly copy the .env file
COPY --from=builder /app/.env .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]