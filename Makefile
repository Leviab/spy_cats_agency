.PHONY: postgres run build test swagger start migrate-up migrate-down migrate-create

postgres:
	@echo "Starting postgres..."
	docker-compose up -d postgres

run:
	@echo "Starting the application..."
	go run ./cmd/api

build:
	@echo "Building the application..."
	go build -o spy_cats_agency ./cmd/api

test:
	@echo "Running tests..."
	go test -v ./...

swagger:
	@echo "Generating swagger docs..."
	swag init -g cmd/api/main.go

migrate-up:
	@echo "Running database migrations..."
	@go run ./cmd/migrate -direction=up

migrate-down:
	@echo "Rolling back database migrations..."
	@go run ./cmd/migrate -direction=down

migrate-create:
	@echo "Creating new migration file..."
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=migration_name"; exit 1; fi
	@migrate create -ext sql -dir db/migration -seq $(name)

start:
	@echo "Starting postgres and waiting for it to be ready..."
	@docker-compose up -d postgres
	@echo "Waiting for postgres to be healthy (this may take 10-15 seconds)..."
	@powershell -Command "Start-Sleep 15"
	@echo "Running database migrations..."
	@go run ./cmd/migrate -direction=up
	@echo "Starting the application..."
	@go run ./cmd/api
