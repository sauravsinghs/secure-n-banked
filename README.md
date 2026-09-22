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
# Simple Bank

Simple Bank is a backend service for a small banking system. It manages users, bank accounts, sessions, email verification, and money transfers with PostgreSQL-backed transactions.

The project is written in Go and exposes the service through gRPC and an HTTP JSON gateway. It also includes database migrations, generated SQL access code, background jobs for email delivery, Docker support, tests, and Kubernetes/AWS deployment manifests.

## What it does

- Creates and authenticates users.
- Issues access and refresh tokens for authenticated sessions.
- Stores user sessions and supports token renewal.
- Sends email verification tasks through a Redis-backed worker queue.
- Verifies user email addresses using verification records.
- Manages bank accounts with supported currencies.
- Performs money transfers between accounts safely using database transactions.
- Tracks account entries and transfer history.
- Runs database migrations automatically on application startup.
- Serves Swagger UI for the HTTP gateway API.

## Tech stack

| Area            | Tools / libraries                                            |
| --------------- | ------------------------------------------------------------ |
| Language        | Go `1.24.4`                                                  |
| HTTP framework  | gRPC-Gateway, Gin package for REST handlers                  |
| RPC             | gRPC, Protocol Buffers                                       |
| Database        | PostgreSQL 12                                                |
| DB access       | `sqlc`, `pgx/v5`                                             |
| Migrations      | `golang-migrate`                                             |
| Authentication  | PASETO tokens, JWT package available                         |
| Background jobs | Redis, `asynq`                                               |
| Email           | Gmail SMTP via `jordan-wright/email`                         |
| Config          | `viper`, `.env` file                                         |
| Logging         | `zerolog`                                                    |
| API docs        | Swagger/OpenAPI generated from protobuf                      |
| Tests/mocks     | Go test, `testify`, `gomock`                                 |
| Containers      | Docker, Docker Compose                                       |
| Deployment      | Kubernetes manifests, AWS/EKS-related config, GitHub Actions |

## Project structure

```text
.
├── api/                 # Gin HTTP handlers, middleware, validators, tests
├── db/
│   ├── migration/       # PostgreSQL schema migrations
│   ├── query/           # SQL queries used by sqlc
│   ├── sqlc/            # Generated database code and transaction logic
│   └── mock/            # Generated database mocks
├── doc/                 # DB docs and generated Swagger UI/assets
├── gapi/                # gRPC server implementation and gateway helpers
├── mail/                # Email sender implementation
├── pb/                  # Generated protobuf Go code
├── proto/               # Protobuf service and message definitions
├── token/               # Token creation and validation
├── util/                # Config, password, random data helpers
├── val/                 # Validation helpers
├── worker/              # Redis/asynq task distributor and processor
├── Dockerfile
├── docker-compose.yaml
├── Makefile
└── main.go              # Application entry point
```

## Prerequisites

Install these tools for local development:

- Go `1.24.4` or compatible newer version
- Docker and Docker Compose
- PostgreSQL client tools, if running the DB manually
- `golang-migrate` CLI, for manual migrations
- `sqlc`, if regenerating database code
- `protoc` and protobuf plugins, if regenerating gRPC/gateway code
- `mockgen`, if regenerating mocks

## Configuration

The application reads configuration from `app.env` using Viper.

Create your local config from the example file:

```sh
cp app.env.example app.env
```

Important variables:

| Variable                 | Purpose                                                 |
| ------------------------ | ------------------------------------------------------- |
| `ENVIRONMENT`            | Runtime environment, for example `development`          |
| `ALLOWED_ORIGINS`        | CORS allowed origins for the HTTP gateway               |
| `POSTGRES_PASSWORD`      | PostgreSQL password used by Docker Compose/start script |
| `DB_SOURCE`              | PostgreSQL connection string                            |
| `MIGRATION_URL`          | Migration location, usually `file://db/migration`       |
| `HTTP_SERVER_ADDRESS`    | HTTP gateway bind address, default `0.0.0.0:8080`       |
| `GRPC_SERVER_ADDRESS`    | gRPC bind address, default `0.0.0.0:9090`               |
| `TOKEN_SYMMETRIC_KEY`    | 32-character symmetric key for PASETO tokens            |
| `ACCESS_TOKEN_DURATION`  | Access token lifetime, for example `15m`                |
| `REFRESH_TOKEN_DURATION` | Refresh token lifetime, for example `24h`               |
| `REDIS_ADDRESS`          | Redis address used by background workers                |
| `EMAIL_SENDER_*`         | Gmail sender settings for verification emails           |

> Do not commit real passwords, token keys, or email credentials.

## Run with Docker Compose

This is the easiest way to start the full stack: PostgreSQL, Redis, and the API server.

```sh
cp app.env.example app.env
# edit app.env and set POSTGRES_PASSWORD, DB_SOURCE, TOKEN_SYMMETRIC_KEY, email values, etc.
docker compose up --build
```

Services exposed by Compose:

- HTTP gateway: `http://localhost:8080`
- gRPC server: `localhost:9090`
- PostgreSQL: `localhost:5432`

The API container waits for PostgreSQL, starts the app, and the app runs migrations automatically.

## Run locally without Docker Compose

Start PostgreSQL and Redis first. You can use the Makefile helpers:

```sh
make postgres
make createdb
make redis
```

Then prepare config:

```sh
cp app.env.example app.env
# edit app.env so DB_SOURCE points to localhost and REDIS_ADDRESS=localhost:6379
```

Run migrations manually if needed:

```sh
make migrateup
```

Start the application:

```sh
make server
```

Equivalent command:

```sh
go run main.go
```

## API usage

The current `main.go` starts:

- a gRPC server on `GRPC_SERVER_ADDRESS` (`:9090` by default)
- an HTTP JSON gateway on `HTTP_SERVER_ADDRESS` (`:8080` by default)
- a Redis/asynq background task processor

### HTTP gateway endpoints

Generated from `proto/service_simple_bank.proto`:

| Method  | Path               | Description                             |
| ------- | ------------------ | --------------------------------------- |
| `POST`  | `/v1/create_user`  | Create a new user                       |
| `POST`  | `/v1/login_user`   | Login and receive access/refresh tokens |
| `PATCH` | `/v1/update_user`  | Update user data                        |
| `GET`   | `/v1/verify_email` | Verify a user's email address           |

Swagger UI is served at:

```text
http://localhost:8080/swagger/
```

### gRPC

The protobuf service is defined in:

```text
proto/service_simple_bank.proto
```

You can connect to the local gRPC server with Evans:

```sh
make evans
```

## Database

The schema includes:

- `users` - registered users, hashed passwords, email verification status, roles
- `accounts` - user-owned bank accounts and balances
- `entries` - account balance change records
- `transfers` - transfer records between accounts
- `sessions` - refresh token sessions and client metadata
- `verify_emails` - email verification codes and expiry data

Migrations are in `db/migration`.

Useful migration commands:

```sh
make migrateup       # apply all migrations
make migrateup1      # apply one migration
make migratedown     # roll back all migrations
make migratedown1    # roll back one migration
make new_migration name=add_something
```

## Development commands

```sh
make test        # run short tests with coverage
make sqlc        # regenerate Go DB code from SQL queries
make mock        # regenerate gomock mocks
make proto       # regenerate protobuf, gRPC gateway, OpenAPI, and statik files
make db_docs     # build DB docs from DBML
make db_schema   # generate SQL schema from DBML
```

## Testing

Run the test suite:

```sh
make test
```

Some database tests expect PostgreSQL to be running with the configured `DB_SOURCE`. Start PostgreSQL and apply migrations before running them locally.

## Docker image

Build the API image:

```sh
docker build -t simplebank .
```

Run it manually after PostgreSQL and Redis are available:

```sh
docker run --env-file app.env -p 8080:8080 -p 9090:9090 simplebank
```

## Deployment notes

The repository contains Kubernetes and AWS/EKS-related files such as:

- `.github/workflows/test.yml`
- `.github/workflows/deploy.yml`
- `eks/`
- `cm.yaml`
- `alb-ingressclass.yaml`
- IAM policy and trust policy JSON files

These files are intended to support CI, container deployment, and cloud infrastructure setup.

## Typical local workflow

```sh
cp app.env.example app.env
# update app.env

docker compose up --build

# in another terminal
make test
```

Then open:

```text
http://localhost:8080/swagger/
```
