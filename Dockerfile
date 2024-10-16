# Dockerfile for the Go storage node

# Step 1: Build the Go app
FROM golang:1.19 AS builder
WORKDIR /app

# Copy and build the application
COPY . .
RUN go mod tidy
RUN go build -o storage-node .

# Step 2: Create a small image for running the Go app
FROM debian:buster
WORKDIR /app

COPY --from=builder /app/storage-node /app/
COPY --from=builder /app/.env /app/.env 
# Add .env if required

# Expose port for communication
EXPOSE 8080

# Start the Go app
CMD ["./storage-node"]
