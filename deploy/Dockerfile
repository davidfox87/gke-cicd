# # Build stage
# FROM golang:1.24.2-alpine AS builder
# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# # Build a fully static binary
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# # Runtime stage
# FROM alpine:latest
# WORKDIR /app
# COPY --from=builder /app/server .
# EXPOSE 8080
# CMD ["./server"]




# Build stage
# FROM golang:1.24.2 AS builder
# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# # Build fully static binary
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# # Runtime stage
# FROM scratch
# WORKDIR /app
# COPY --from=builder /app/server .
# EXPOSE 8080
# CMD ["./server"]


# -----------------------------
# Build stage
# -----------------------------
FROM golang:1.24.2 AS builder

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build fully static binary
# - CGO_ENABLED=0 ensures no dynamic linking
# - GOOS=linux GOARCH=amd64 target Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# -----------------------------
# Runtime stage
# -----------------------------
FROM scratch
WORKDIR /app

# Copy the statically built binary from the builder
COPY --from=builder /app/server .

# Expose port
EXPOSE 8080

# Start the server
CMD ["./server"]
    