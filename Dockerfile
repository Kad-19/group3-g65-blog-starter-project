# --- STAGE 1: Build ---
# Use the official Go image. 'alpine' is a lightweight Linux distribution.
FROM golang:1.24.6-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files.
COPY go.mod go.sum ./

# Download external dependencies.
RUN go mod download

# Copy the entire source code into the container.
COPY . .

# ====================================================================
# === VITAL STEP: Tidy modules to recognize all local packages ===
RUN go mod tidy
# ====================================================================

# Build the Go application.
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/api

# --- STAGE 2: Final/Runtime ---
# Use a minimal, secure base image.
FROM alpine:latest

# Add certificates for HTTPS requests.
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy ONLY the compiled binary from the 'builder' stage.
COPY --from=builder /server .

# Copy templates folder
COPY utils ./utils

# Expose port 8080
EXPOSE 8080

# This is the command that will run when the container starts.
CMD ["./server"]