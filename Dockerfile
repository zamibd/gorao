# Build stage
FROM golang:1.25.7-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# We use the Makefile's build logic but adapt it for Docker
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o gorao .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates curl tzdata libcap

# Set timezone
ENV TZ=Asia/Dhaka

# Create a non-root user
RUN adduser -D -g '' -u 1000 gorao

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/gorao .

# Change ownership of the binary and allow it to bind to privileged ports
RUN chown gorao:gorao /app/gorao && \
    setcap 'cap_net_bind_service=+ep' /app/gorao

# Expose ports
# DNS
EXPOSE 53/udp 53/tcp
# HTTP/HTTPS
EXPOSE 80/tcp 443/tcp 443/udp
# HTTPS (SNI TLS)
EXPOSE 8443/tcp
# DoT
EXPOSE 853/tcp
# DoH
EXPOSE 8053/tcp
# DoQ
EXPOSE 8853/udp

# Healthcheck
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f --insecure "https://127.0.0.1:8053/dns-query?dns=EjQBAAABAAAAAAAAA3d3dwZnb29nbGUDY29tAAABAAE" || exit 1

# Switch to non-root user
USER gorao

# Entrypoint
ENTRYPOINT ["/app/gorao"]
