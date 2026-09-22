# Secure N Banked

Secure N Banked is a backend service for a simple banking system. It provides user registration, authentication, email verification, account management, money transfers, refresh-token sessions, and background email tasks.

The service exposes both a gRPC API and a REST/JSON API through grpc-gateway. Database schema changes are managed with versioned migrations, and generated SQL access code is maintained with sqlc.

## Tech Stack

- Go 1.24.4
- PostgreSQL 12
- Redis 8
- gRPC and Protocol Buffers
- grpc-gateway for REST/JSON endpoints
- pgx and sqlc for PostgreSQL access
- golang-migrate for database migrations
- Asynq for background task processing
- JWT and PASETO token support
- Docker and Docker Compose

## API Endpoints

The REST gateway runs on `http://localhost:8080` and includes:

- `POST /v1/create_user`
- `PATCH /v1/update_user`
- `POST /v1/login_user`
- `GET /v1/verify_email`

The gRPC server runs on `localhost:9090`.

## Prerequisites

- Go 1.24 or newer
- Docker Desktop with Docker Compose
- Git

## Configuration

The application reads configuration from `app.env`. A template is available in `app.env.example`.

The default local setup uses:

- PostgreSQL: `localhost:5433`
- Redis: `localhost:6379`
- HTTP gateway: `localhost:8080`
- gRPC server: `localhost:9090`

Port `5433` is used for Docker PostgreSQL so it can coexist with a PostgreSQL installation already using port `5432` on Windows.

## Run Locally

From the repository root, start PostgreSQL and Redis:

```powershell
docker compose up -d --wait postgres redis
```

Start the Go application:

```powershell
go run .
```

The application runs database migrations automatically during startup.

Stop the application with `Ctrl+C`. Stop the Docker dependencies with:

```powershell
docker compose down
```

Do not use `docker compose down -v` unless you want to delete the database volume and its data.

## Run Tests

Start the database and Redis if they are not already running, then run all tests:

```powershell
docker compose up -d --wait postgres redis
go test ./...
```

The SQL integration tests run their required migrations automatically.

## Useful Commands

```powershell
# Run tests with coverage output
go test -v -short -cover ./...

# Run the server
go run .

# Apply migrations manually
migrate -path db/migration -database "postgres://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up

# Generate sqlc code
sqlc generate
```

## Run Everything with Compose

To build and run the API together with PostgreSQL and Redis:

```powershell
docker compose up --build
```

The containerized API uses PostgreSQL at `postgres:5432` internally and is published on ports `8080` and `9090`.
