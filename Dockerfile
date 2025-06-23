# Stage 1: Build the application using a Debian-based Go image
FROM golang:1.23-bookworm AS builder

# Install build dependencies for SQLite and plugins using apt-get
# build-essential includes gcc, binutils (for ld), and other necessary tools
RUN apt-get update && apt-get install -y build-essential libsqlite3-dev

WORKDIR /src

# Copy the entire backend source code
COPY casino-backend/ ./

# Build the game plugins
WORKDIR /src/games/example
RUN go mod download
RUN go build -buildmode=plugin -o /src/games/example.so .

# Build the main application
WORKDIR /src/casino
RUN go mod tidy
# CGO_ENABLED is required for the sqlite driver
RUN CGO_ENABLED=1 go build -o /casino-app main.go

# Stage 2: Create the final image using a minimal Debian base
FROM debian:bookworm-slim

# Install runtime dependencies for SQLite
RUN apt-get update && apt-get install -y libsqlite3-0 && rm -rf /var/lib/apt/lists/*

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
