# Build stage
FROM golang:1.24.2 AS builder 

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN go build -o server .

# Runtime stage
FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]