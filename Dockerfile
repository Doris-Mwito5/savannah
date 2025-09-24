# Multi-stage build
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /usr/bin
COPY --from=builder /app/main .
COPY --from=builder /app/.env .env

EXPOSE 8080
CMD ["./main"]