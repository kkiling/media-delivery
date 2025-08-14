# Build stage
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN go build -o server_app ./cmd/main.go

# Runtime stage
FROM lscr.io/linuxserver/baseimage-ubuntu:jammy

# Install runtime deps
RUN apt-get update && \
    apt-get install -y mkvtoolnix && \
    rm -rf /var/lib/apt/lists/*

# Copy binary
COPY --from=builder /app/server_app /app/server_app

# Set working directory
WORKDIR /app

# Expose ports
EXPOSE 8080

# Run as non-root user
USER abc

CMD ["/app/server_app", "-c", "/config/config.yaml"]