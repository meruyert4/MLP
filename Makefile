.PHONY: run build migrate docker-up docker-down docker-logs clean test

run:
	go run cmd/main.go

build:
	go build -o bin/mlp cmd/main.go

migrate:
	@echo "Running migrations..."
	@if [ -f .env ]; then \
		export $$(cat .env | grep -v '^#' | xargs) && \
		psql $${DB_HOST:-localhost}:$${DB_PORT:-5432}/$${DB_NAME:-mlp} -U $${DB_USER:-postgres} -f migrations/001_init.sql; \
	else \
		echo "Error: .env file not found"; \
		exit 1; \
	fi
	@echo "Migrations completed"

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-build:
	docker-compose build

docker-restart:
	docker-compose restart

clean:
	rm -rf bin/
	docker-compose down -v

test:
	go test -v ./...

deps:
	go mod download
	go mod tidy

fmt:
	go fmt ./...

lint:
	golangci-lint run

help:
	@echo "Available commands:"
	@echo "  make run           - Run the application locally"
	@echo "  make build         - Build the application binary"
	@echo "  make migrate       - Run database migrations"
	@echo "  make docker-up     - Start all services with docker-compose"
	@echo "  make docker-down   - Stop all services"
	@echo "  make docker-logs   - Show docker logs"
	@echo "  make docker-build  - Build docker images"
	@echo "  make docker-restart- Restart docker services"
	@echo "  make clean         - Clean build artifacts and stop docker"
	@echo "  make test          - Run tests"
	@echo "  make deps          - Download and tidy dependencies"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Run linter"
