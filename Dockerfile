# --- STAGE 1: Build ---
# Use the official Go image. 'alpine' is a lightweight Linux distribution.
FROM golang:1.24.6-alpine AS builder
WORKDIR /app
LABEL maintainer="Temesgen <guys199421@gmail.com>"
LABEL version="1.0"
LABEL description="Group3 G65 Blog Starter Project API"
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# --- STAGE 2: Final/Runtime ---
FROM alpine:3.20
RUN apk --no-cache add ca-certificates
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /home/appuser/
COPY --from=builder /app/server .
COPY utils ./utils
ENV GIN_MODE=release
ENV PORT=8080
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
	CMD wget --spider --quiet http://localhost:8080/health || exit 1
CMD ["./server", "--host", "0.0.0.0", "--port", "8080"]