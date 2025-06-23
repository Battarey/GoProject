FROM golang:1.23-alpine as builder
WORKDIR /app
COPY . .
RUN go build -o api-gateway main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-gateway .
EXPOSE 8080
CMD ["./api-gateway"]
