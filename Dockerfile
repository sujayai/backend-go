FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install git for Go modules
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Create data directories
RUN mkdir -p data posts uploads

# Expose port
EXPOSE 3000

# Run the server
CMD ["./server"]
