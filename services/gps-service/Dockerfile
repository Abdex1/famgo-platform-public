#syntax=docker/dockerfile:1

# === Build stage: Compile Go application ===
FROM dhi.io/golang:1-alpine3.22-dev AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gps-service ./cmd

# === Final stage: Create minimal runtime image ===
FROM dhi.io/alpine-base:3.22

WORKDIR /app

COPY --from=builder /build/gps-service .

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD ["wget", "--quiet", "--tries=1", "--spider", "http://localhost:5002/health", "||" , "exit", "1"]

EXPOSE 5002

CMD ["./gps-service"]
