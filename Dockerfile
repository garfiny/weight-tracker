# ---- Build stage ----
FROM golang:1.24-4 AS builder

WORKDIR /src

# Set target platform and disable CGO for static binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Copy go.mod and go.sum first for dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy only necessary source files (avoid .git, docs, etc.)
COPY . .

# Build the binary with optimizations
RUN go build -trimpath -ldflags="-s -w" -o /out/weight-tracker ./src

# ---- Runtime stage ----
FROM alpine:3.20

# Add non-root user and minimal runtime deps
RUN addgroup -S app && adduser -S app -G app

# Create app directory
WORKDIR /app

# Copy only what’s needed
COPY --from=builder /out/weight-tracker /usr/local/bin/
COPY --from=builder /src/webapp /app/webapp

# Ensure ownership and drop privileges
RUN chown -R app:app /usr/local/bin/weight-tracker /app/webapp
USER app

# Minimal runtime environment
EXPOSE 8080
ENV PORT=8080

ENTRYPOINT ["/usr/local/bin/weight-tracker"]
