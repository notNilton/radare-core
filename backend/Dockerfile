# Stage 1: Build the Go application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]