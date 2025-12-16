.PHONY: build run test migrate-up migrate-down migrate-create clean

# Build the application
build:
	@echo "Building StudySync..."
	go build -o bin/studysync ./cmd/main.go

# Run the application
run:
	@echo "Running StudySync..."
	go run ./cmd/main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test ./internal/api -v -tags=integration

# Run database migrations up
migrate-up:
	@echo "Running migrations up..."
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/studysync?sslmode=disable" up

# Run database migrations down
migrate-down:
	@echo "Running migrations down..."
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/studysync?sslmode=disable" down 1

# Create new migration
migrate-create:
	@echo "Creating new migration..."
	migrate create -ext sql -dir migrations -seq $(name)

# Clean binaries
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -rf logs/

# Generate Swagger docs
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/main.go -o docs

# Docker compose up
docker-up:
	docker-compose up -d

# Docker compose down
docker-down:
	docker-compose down

# View logs
logs:
	docker-compose logs -f api

# Help
help:
	@echo "Available commands:"
	@echo "  build               - Build the application"
	@echo "  run                 - Run the application"
	@echo "  test                - Run all tests"
	@echo "  test-integration    - Run integration tests"
	@echo "  migrate-up          - Run database migrations up"
	@echo "  migrate-down        - Rollback last migration"
	@echo "  migrate-create name - Create new migration file"
	@echo "  clean               - Clean binaries and logs"
	@echo "  swagger             - Generate Swagger docs"
	@echo "  docker-up           - Start Docker containers"
	@echo "  docker-down         - Stop Docker containers"
	@echo "  logs                - View container logs"