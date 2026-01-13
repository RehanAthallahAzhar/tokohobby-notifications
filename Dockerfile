# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy messaging library (local dependency)
COPY messaging /app/messaging

# Copy notifications service
COPY notifications /app/notifications

WORKDIR /app/notifications

# Download dependencies
RUN go mod download

# Build notification worker
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notification-worker ./cmd/worker

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/notifications/notification-worker .

# Run notification worker
CMD ["./notification-worker"]
