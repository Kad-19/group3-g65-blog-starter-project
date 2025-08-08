# --- STAGE 1: Build ---
# Use the official Go image. 'alpine' is a lightweight Linux distribution.
FROM golang:1.24.4-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files. This is cached by Docker and will only
# re-run if go.mod or go.sum changes, speeding up future builds.
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code into the container.
# We need all packages (delivery, repository, etc.) to build the app.
COPY . .

# Build the Go application.
# CGO_ENABLED=0 creates a statically-linked binary.
# The final argument './cmd/api' tells the compiler where your main package is.
# The output will be a single file named 'server' in the root directory.
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/api

# --- STAGE 2: Final/Runtime ---
# Use a minimal, secure base image. Alpine is very small.
FROM alpine:latest

# Alpine doesn't have root CA certificates by default, which are needed
# to make HTTPS calls to external services (like AI APIs, email servers, etc.).
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy ONLY the compiled binary from the 'builder' stage.
# This is the magic of multi-stage builds. The final image will be tiny
# and won't contain any source code or build tools.
COPY --from=builder /server .

# Expose port 8080 to the outside world. This is what your app listens on.
EXPOSE 8080

# This is the command that will run when the container starts.
# It simply executes your compiled application.
CMD ["./server"]