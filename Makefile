DB_URL=postgres://user:password@localhost:5432/spy_cats_agency?sslmode=disable

.PHONY: postgres createdb dropdb migrateup migratedown run build test swagger

postgres:
	@echo "Starting postgres..."
	docker-compose up -d postgres

createdb:
	@echo "Creating database..."
	docker-compose exec -T postgres createdb --username=user --owner=user spy_cats_agency


dropdb:
	@echo "Dropping database..."
	docker-compose exec -T postgres dropdb spy_cats_agency

migrateup:
	@echo "Running migrations up..."
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	@echo "Running migrations down..."
	migrate -path db/migration -database "$(DB_URL)" -verbose down

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

setup:
	@make postgres
	@echo "Waiting for postgres to be healthy..."
	@sleep 5
	@-make createdb
	@make migrateup

start: setup
	@make run
