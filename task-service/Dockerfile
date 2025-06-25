FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o task-service main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/task-service .
CMD ["./task-service"]
