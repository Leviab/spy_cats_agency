# Spy Cat Agency - Management System

A RESTful API for managing spy cats, missions, and targets built with Go, Gin, and PostgreSQL.

## Features

- **Spy Cat Management**: Create, read, update, and delete spy cats with breed validation via TheCatAPI
- **Mission Management**: Create missions with 1-3 targets, assign cats, and track completion
- **Target Management**: Update notes, mark targets as complete, and manage target lifecycle
- **Business Rules**: Enforces all specified constraints (one mission per cat, target limits, completion rules)
- **API Documentation**: Auto-generated Swagger/OpenAPI documentation

## Quick Start

### Prerequisites

- Go 1.19 or later
- Docker and Docker Compose
- Make (optional, for convenience commands)

### Setup and Run

1. **Clone and setup the application:**
   ```bash
   git clone https://github.com/Leviab/spy_cats_agency
   cd spy_cats_agency
   cp .env.example .env
   ```

2. **Start the application:**
   ```bash
   make start
   ```

The application will:
- Start PostgreSQL in Docker
- Run database migrations
- Start the API server on port 8080

### Alternative Setup (Manual)

If you prefer manual setup:

```bash
# Start database
docker-compose up -d postgres

# Wait for database to be ready, then run migrations
migrate -path db/migration -database "postgres://user:password@localhost:5432/spy_cats_agency?sslmode=disable" up

# Start the application
go run ./cmd/api
```

## API Documentation

- **Swagger UI**: http://localhost:8080/swagger/index.html

## Development

### Available Make Commands

```bash
make postgres    # Start PostgreSQL container
make migrateup   # Run database migrations
make migratedown # Rollback database migrations
make run         # Start the application
make build       # Build the application binary
make test        # Run tests
make swagger     # Generate Swagger documentation
make setup       # Setup database and run migrations
make start       # Complete setup and start application
```

## Configuration

The application uses environment variables for configuration. Copy `.env.example` to `.env` and modify as needed:

```env
DB_HOST=localhost
DB_USER=user
DB_PASSWORD=password
DB_NAME=spy_cats_agency
DB_PORT=5432
SERVER_PORT=8080
CAT_API_ENDPOINT=https://api.thecatapi.com/v1
```

## Testing

Run the test suite:

```bash
make test
```