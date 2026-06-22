.PHONY: install build start dev test test-watch test-cov lint format clean docker-build docker-run db-migrate db-generate help

# Default target
help:
	@echo "FamGo Service Template - Available Commands"
	@echo "==========================================="
	@echo "make install            - Install dependencies"
	@echo "make build              - Build production code"
	@echo "make start              - Start production server"
	@echo "make dev                - Start development server with hot reload"
	@echo "make start:debug        - Start with debugger"
	@echo "make test               - Run unit tests"
	@echo "make test:watch         - Run tests in watch mode"
	@echo "make test:cov           - Run tests with coverage"
	@echo "make lint               - Run ESLint"
	@echo "make format             - Format code with Prettier"
	@echo "make clean              - Remove dist/ node_modules/"
	@echo "make db-generate        - Generate database migration"
	@echo "make db-migrate         - Run database migrations"
	@echo "make db-revert          - Revert last migration"
	@echo "make docker-build       - Build Docker image"
	@echo "make docker-run         - Run Docker container"

install:
	npm install

build:
	npm run build

start:
	npm start

dev:
	npm run start:dev

start:debug:
	npm run start:debug

test:
	npm test

test:watch:
	npm run test:watch

test:cov:
	npm run test:cov

lint:
	npm run lint

format:
	npm run format

clean:
	rm -rf dist node_modules coverage .nyc_output

db-generate:
	npm run migration:generate

db-migrate:
	npm run migration:run

db-revert:
	npm run migration:revert

docker-build:
	docker build -t famgo/template-service:latest .

docker-run:
	docker run -p 3000:3000 \
		-e NODE_ENV=development \
		-e DB_HOST=host.docker.internal \
		-e REDIS_HOST=host.docker.internal \
		famgo/template-service:latest

.dockerignore-generate:
	@echo "node_modules" > .dockerignore
	@echo "dist" >> .dockerignore
	@echo ".git" >> .dockerignore
	@echo ".gitignore" >> .dockerignore
	@echo "*.md" >> .dockerignore
	@echo ".env" >> .dockerignore
	@echo ".env.local" >> .dockerignore
