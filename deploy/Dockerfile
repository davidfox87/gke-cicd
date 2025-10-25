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
COPY *.go ./

# Build fully static binary
# - CGO_ENABLED=0 ensures no dynamic linking
# - GOOS=linux GOARCH=amd64 target Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /server

# -----------------------------
# Runtime stage
# -----------------------------
# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /

# Copy the statically built binary from the builder
COPY --from=builder /server /server

# Expose port
EXPOSE 8080

USER nonroot:nonroot

# Start the server
CMD ["./server"]
    