# # -----------------------------
# # Build stage
# # -----------------------------
# FROM golang:1.25 AS builder

# # Set working directory inside the container
# WORKDIR /app

# # Copy go.mod and go.sum first to leverage caching
# COPY go.mod go.sum ./
# RUN go mod download

# # Copy the rest of the source code
# COPY . .

# # Build fully static binary
# # - CGO_ENABLED=0 ensures no dynamic linking
# # - GOOS=linux GOARCH=amd64 target Linux
# RUN CGO_ENABLED=0 GOOS=linux go build -o /server

# # -----------------------------
# # Runtime stage
# # -----------------------------
# # Deploy the application binary into a lean image
# FROM gcr.io/distroless/base-debian11 AS build-release-stage
# WORKDIR /

# # Copy the statically built binary from the builder
# COPY --from=builder /server /server

# # Expose port
# EXPOSE 8080

# USER nonroot:nonroot

# # Start the server
# CMD ["./server"]
# -----------------------------
# Build stage
# -----------------------------
FROM golang:1.25 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Install Delve statically and verify
RUN CGO_ENABLED=0 GOBIN=/go/bin go install github.com/go-delve/delve/cmd/dlv@v1.23.0 && \
    ls -l /go/bin/dlv

RUN CGO_ENABLED=0 GOOS=linux go build -gcflags "all=-N -l" -o /server

# -----------------------------
# Runtime stage (distroless)
# -----------------------------
FROM scratch

# No WORKDIR needed, but can add for clarity
COPY --from=builder /server /server
COPY --from=builder /go/bin/dlv /dlv

# No EXPOSE needed for port-forward, but good for docs
EXPOSE 8080 2345

# Start Delve as DAP server
# --listen=:2345 binds immediately
# --log for startup debug (remove later)
# --accept-multiclient not needed since dap mode is single-client by default, but add if you want
ENTRYPOINT ["/dlv", "dap", "--listen=0.0.0.0:2345", "--log"]
