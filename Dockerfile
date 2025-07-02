FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

ARG GITHUB_TOKEN
ENV GOPRIVATE=github.com/k5sha/*

RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o chat-service ./cmd/chat-service

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/chat-service .
COPY --from=builder /app/cmd/migrate/migrations ./cmd/migrate/migrations

CMD ["./chat-service"]
