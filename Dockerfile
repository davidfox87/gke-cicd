FROM golang:1.25 AS builder

ARG OS="linux"
ARG ARCH="amd64"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1 GOOS=${OS} GOARCH=${ARCH}

RUN go build -o /bin/api ./cmd/api

FROM gcr.io/distroless/base-debian12

COPY --from=builder /bin/api /api

EXPOSE 8080

ENTRYPOINT ["/api"]
