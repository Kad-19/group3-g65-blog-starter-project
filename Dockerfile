# --- STAGE 1: Build the application binary ---
# Use the official Go image as a builder.
# 'alpine' is a lightweight version of Linux, which keeps the build stage small.
FROM golang:1.24.4-alpine AS builder

# Set the working directory inside the container for the build
WORKDIR /app

# Copy the Go module files. This is a key optimization.
# Docker caches this layer. Dependencies are only re-downloaded if go.mod or go.sum changes.
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the Go application.
# -o /app/main specifies the output path and name for our compiled binary.
# The build command targets your specific entry point: cmd/api/main.go.
# CGO_ENABLED=0 creates a statically linked binary, which is highly portable.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api/main.go


# --- STAGE 2: Create the final, lightweight runtime image ---
# Use a minimal, secure base image. 'alpine' is small and has a package manager.
# 'scratch' is even smaller but has nothing, not even CA certificates for HTTPS calls.
# Alpine is a safer starting point.
FROM alpine:3.20

WORKDIR /app

# We need root certificates for making external HTTPS requests (e.g., to AI services, email servers).
RUN apk --no-cache add ca-certificates

# Copy ONLY the compiled binary from the 'builder' stage.
# This is the magic of multi-stage builds! Our final image is tiny and secure
# because it doesn't contain the source code or Go compiler.
COPY --from=builder /app/main .

# Copy your .env file which contains all the configuration.
# This ensures your app can find its configuration when it starts.
COPY .env .

# Expose the port that your Go application listens on.
EXPOSE 8080

# This is the command that will be run when the container starts.
# It simply executes your compiled application binary.
CMD ["./main"]