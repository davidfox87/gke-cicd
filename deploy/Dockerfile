FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

# Build the Go binary
RUN go build -o server .

# Use a minimal base image for runtime
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080

CMD ["./server"]
