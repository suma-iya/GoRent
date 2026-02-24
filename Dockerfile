# Build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

# Final stage
FROM debian:bookworm-slim

# Install CA certificates for HTTPS requests
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary
COPY --from=builder /app/main .

# Copy config files (including Firebase service account if needed)
COPY config/ ./config/

EXPOSE 8080

CMD ["./main"]
