# Step 1: Build the Go binary
FROM golang:alpine AS builder

# Install git for resolving Go module dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# We build from the cmd/server directory where main.go is located
RUN CGO_ENABLED=0 GOOS=linux go build -o dexra-backend ./cmd/server

# Step 2: Create a lightweight runtime image
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS requests, and compat libraries for ONNX runtime (Chroma Go)
RUN apk --no-cache add ca-certificates libstdc++ libc6-compat gcompat

# Copy the production environment file
COPY .env.prod .env

# Copy the binary from the builder stage
COPY --from=builder /app/dexra-backend .

# Expose the port the app runs on (update if your server uses a different port)
EXPOSE 9008

# Command to run the executable
CMD ["./dexra-backend"]
