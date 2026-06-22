FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o auth-service ./cmd/api

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/auth-service .

EXPOSE 8080

CMD ["./auth-service"]