# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

# Install build dependencies for SQLite and plugins.
# Explicitly add binutils to ensure the linker (ld) is available,
# which is required for building Go plugins.
RUN apk add --no-cache gcc musl-dev sqlite binutils

WORKDIR /src

# Copy the entire backend source code
# This is necessary because of the 'replace' directive in go.mod
COPY casino-backend/ ./

# Build the game plugins
# First, download dependencies for the plugin
WORKDIR /src/games/example
RUN go mod download
# Then, build the plugin
RUN go build -buildmode=plugin -o /src/games/example.so .

# Build the main application
WORKDIR /src/casino
# Tidy up dependencies
RUN go mod tidy
# Build the binary
# CGO_ENABLED is required for the sqlite driver
RUN CGO_ENABLED=1 go build -o /casino-app main.go

# Stage 2: Create the final image
FROM alpine:latest

# Install runtime dependencies for SQLite
RUN apk add --no-cache sqlite

# Set environment to production so the correct DB path is used
ENV ENV=production

WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /casino-app .

# Create a directory for game plugins and copy them from the builder stage
# The application looks for plugins in ../games, which resolves to /games from /app
RUN mkdir /games
COPY --from=builder /src/games/example.so /games/

# Expose the port the server listens on
EXPOSE 9000

# Create a volume for the persistent database
VOLUME ["/data"]

# Command to run the application
CMD ["./casino-app"]
