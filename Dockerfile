# --- STAGE 1: Build ---
# Use the official Go image. 'alpine' is a lightweight Linux distribution.
FROM golang:1.24.6-alpine3.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Metadata labels (for maintainability)
LABEL maintainer="Temesgen <guys199421@gmail.com>"
LABEL version="1.0"
LABEL description="Group3 G65 Blog Starter Project API"

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/api

# --- STAGE 2: Final/Runtime ---
# Use a minimal, secure base image.
FROM alpine:3.20

# Add certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /home/appuser/

# Copy ONLY the compiled binary from the 'builder' stage
COPY --from=builder /server .

# Copy templates folder
COPY utils ./utils

# Set environment variables (example)
ENV GIN_MODE=release
ENV PORT=8080

# Change to non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Healthcheck for container monitoring
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
	CMD wget --spider --quiet http://localhost:8080/health || exit 1

# Start the server
CMD ["./server", "--host", "0.0.0.0", "--port", "8080"]