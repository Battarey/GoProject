# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates git protobuf

# Установка protoc-gen-go и protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH="/root/go/bin:${PATH}"

COPY go.mod go.sum ./
COPY proto ./proto
RUN go clean -modcache && go mod download

COPY . .

# Генерация gRPC файлов
RUN protoc --proto_path=./proto --go_out=./proto --go-grpc_out=./proto ./proto/user.proto

RUN go build -o user-service ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user-service .

EXPOSE 50051

CMD ["./user-service"]
