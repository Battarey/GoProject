# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go clean -modcache && go mod download

COPY . .

RUN go build -o user-service ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user-service .
COPY wait-for-it.sh .
RUN chmod +x wait-for-it.sh

EXPOSE 50051

CMD ["./user-service"]
