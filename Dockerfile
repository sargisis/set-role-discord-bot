# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy source code
COPY . .

# Download dependencies
RUN go mod download

# Build binary
RUN go build -o main main.go

# Run Stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port (Render sets PORT env var)
EXPOSE 8080

# Run
CMD ["./main"]
