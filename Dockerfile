# Use a minimal Go base image
FROM golang:1.23-alpine

# Set environment variable for SQLite DB location (can be overridden)
ENV DB_PATH=/data/db.sqlite

# Install SQLite driver dependencies (for CGO)
RUN apk add --no-cache gcc musl-dev sqlite

# Create app directory
WORKDIR /app

# Copy Go modules and download dependencies
COPY casino-backend/casino/go.mod casino-backend/casino/go.sum ./
RUN go mod download

# Copy the source code
COPY casino-backend/casino/ ./

# Build the binary
RUN go build -o casino-app main.go

# Expose port if needed (e.g., 8080)
EXPOSE 9000

# Create volume directory for persistent DB
RUN mkdir -p /data
VOLUME ["/data"]

# Run the app
CMD ["./casino-app"]
