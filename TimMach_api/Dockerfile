# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.22.1
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

FROM debian:12-slim

WORKDIR /app
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates tzdata fonts-dejavu wkhtmltopdf && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY db/migrations ./db/migrations
COPY entrypoint.sh /app/entrypoint.sh
COPY templates ./templates

RUN chmod +x /app/entrypoint.sh

ENV PORT=8080

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/server"]
