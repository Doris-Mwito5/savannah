FROM golang:1.24-alpine AS builder 

WORKDIR /app
COPY go.* ./
RUN go mod download

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata bash

WORKDIR /root/

# Copy app binary
COPY --from=builder /app/main .
# Copy goose binary
COPY --from=builder /go/bin/goose /usr/local/bin/goose
# Copy migrations
COPY --from=builder /app/internal/db/migrations /migrations

EXPOSE 8080
CMD ["./main"]
