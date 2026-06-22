.PHONY: help install build start dev test test-watch test-cov lint format clean docker-build docker-run db-migrate db-generate

# Default target
help:
	@echo "FamGo Go Service Template - Available Commands"
	@echo "==========================================="
	@echo "make install            - Download dependencies"
	@echo "make build              - Build binary"
	@echo "make start              - Start production server"
	@echo "make dev                - Start development server"
	@echo "make test               - Run unit tests"
	@echo "make test:watch         - Run tests in watch mode"
	@echo "make test:cov           - Run tests with coverage"
	@echo "make lint               - Run Go linter"
	@echo "make format             - Format code"
	@echo "make clean              - Remove bin/ and clean cache"
	@echo "make docker-build       - Build Docker image"
	@echo "make docker-run         - Run Docker container"
	@echo "make db-migrate         - Run database migrations"

install:
	go mod download
	go mod tidy

build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/service cmd/service/main.go

start: build
	./bin/service

dev:
	go run cmd/service/main.go

test:
	go test -v ./...

test:watch:
	go test -v ./... -watch

test:cov:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint:
	golangci-lint run

format:
	go fmt ./...
	goimports -w .

clean:
	rm -rf bin/
	go clean -cache
	rm -f coverage.out

docker-build:
	docker build -t famgo/template-service:latest .

docker-run:
	docker run -p 5001:5001 -p 8080:8080 \
		-e SERVICE_NAME=template-service \
		-e DB_HOST=host.docker.internal \
		-e REDIS_HOST=host.docker.internal \
		famgo/template-service:latest

db-migrate:
	migrate -path migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME" up

db-rollback:
	migrate -path migrations -database "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME" down
