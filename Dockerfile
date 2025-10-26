# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy everything including vendored dependencies
COPY . .

# Build using vendored dependencies
# Note: vendor directory is committed to repository to enable builds in
# environments with restricted network access or TLS certificate issues
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o notes-server .

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/notes-server .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./notes-server"]
