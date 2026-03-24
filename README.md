# Go REST API Template

Go REST API Template is a Golang backend template for building RESTful APIs with a clean, scalable project structure. It is designed as a starting point for Go services that need clear separation of concerns, transaction handling, reusable infrastructure components, and predictable HTTP application flow.

This template includes:

- Clean architecture style layering across `controller`, `service`, `domain`, and `infra`
- A working REST API example for `users`
- PostgreSQL integration with `pgx`
- Transaction manager patterns for multi-step business operations
- Shared HTTP response and error handling helpers
- Environment-based configuration loading
- A layout that is easy to extend for new modules and services

## Use Case

Use this repository as a foundation when you want to start a new Go backend service without rebuilding the same boilerplate each time. It is suited for:

- Internal APIs
- CRUD services
- Modular monolith backends
- Services that need transaction-aware business logic
- Teams standardizing how Go services are organized

## Tech Stack

- Go `1.25.7`
- `chi` for HTTP routing
- `cors` middleware
- `pgx/v5` for PostgreSQL access
- `envconfig` for environment configuration
- `slog` for structured logging

## Architecture

The codebase follows a practical clean architecture approach:

- `internal/controller` handles HTTP transport, request decoding, response mapping, and route registration
- `internal/service` contains application use cases and business orchestration
- `internal/domain` defines core entities, domain errors, and transaction options
- `internal/infra` contains infrastructure implementations such as PostgreSQL repositories and transaction management
- `pkg` contains reusable shared packages such as DTOs and database connection helpers
- `cmd/api` wires everything together and starts the HTTP server

Request flow:

`HTTP request -> controller -> service -> repository/transaction manager -> PostgreSQL`

This keeps transport logic, business logic, and persistence concerns isolated from each other.

## Project Structure

```text
cmd/api/                         Application entrypoint
internal/config/                Environment configuration
internal/controller/            HTTP router and handlers
internal/controller/user/       Example user REST module
internal/domain/                Core entities, errors, transaction options
internal/service/user/          User use cases and service interfaces
internal/infra/postgres/        PostgreSQL abstractions and implementations
internal/infra/postgres/user/   User repository implementation
internal/infra/postgres/tx-manager/ Transaction manager implementation
internal/utils/httputil/        Shared HTTP helpers
pkg/dto/                        Request and response DTOs
pkg/pgxpool/                    PostgreSQL pool bootstrap helpers
```

## Template Features

### 1. REST API module example

The `user` module demonstrates how to organize:

- route registration
- request DTOs
- handler-to-service mapping
- response serialization
- validation entry points

Implemented routes under `/api/v1/users`:

- `GET /`
- `POST /`
- `GET /{id}`
- `PUT /{id}`
- `DELETE /{id}`
- `POST /upsert`

### 2. Transaction handling pattern

The template includes a transaction manager in `internal/infra/postgres/tx-manager`.

Business logic can wrap multi-step workflows with:

```go
err := txm.Do(ctx, func(ctx context.Context) error {
    // use repositories with the transactional context
    return nil
})
```

Repositories automatically pick the active transaction from context through a shared query engine abstraction. This allows service methods to remain explicit about transactional boundaries while keeping repository code simple.

### 3. Shared domain error mapping

The template defines domain-level errors such as:

- `ErrNotFound`
- `ErrInvalidInput`
- `ErrUnauthorized`
- `ErrForbidden`
- `ErrAlreadyExists`

Handlers use a shared HTTP helper to consistently translate domain errors into HTTP status codes and JSON error responses.

### 4. Config bootstrap

Configuration is loaded from environment variables, making the template easy to run locally and in containerized environments.

Supported variables:

- `APP_CONFIG_PORT` default: `8080`
- `LOG_LEVEL` default: `info`
- `DATABASE_URL`
- `APP_CONFIG_POSTGRES_HOST` default: `localhost`
- `APP_CONFIG_POSTGRES_PORT` default: `5432`
- `APP_CONFIG_POSTGRES_USERNAME`
- `APP_CONFIG_POSTGRES_PASSWORD`
- `APP_CONFIG_POSTGRES_DBNAME`

If `DATABASE_URL` is set, it is used directly. Otherwise, the PostgreSQL DSN is built from the individual Postgres variables.

## Getting Started

### 1. Install dependencies

```bash
go mod download
```

### 2. Set environment variables

Example:

```bash
export APP_CONFIG_PORT=8080
export APP_CONFIG_POSTGRES_HOST=localhost
export APP_CONFIG_POSTGRES_PORT=5432
export APP_CONFIG_POSTGRES_USERNAME=postgres
export APP_CONFIG_POSTGRES_PASSWORD=postgres
export APP_CONFIG_POSTGRES_DBNAME=app
```

Or use:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/app?sslmode=disable"
```

### 3. Prepare the database

Create a `users` table compatible with the example module:

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### 4. Run the API

```bash
go run ./cmd/api
```

The server starts on `http://localhost:8080` unless overridden by configuration.

## Example Requests

Create a user:

```bash
curl -X POST http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'
```

List users:

```bash
curl -X GET http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{"limit":10,"offset":0}'
```

Upsert a user:

```bash
curl -X POST http://localhost:8080/api/v1/users/upsert \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'
```

## How To Extend This Template

To add a new module such as `orders`, `products`, or `accounts`:

1. Add domain models and domain-specific errors if needed.
2. Create service interfaces and business logic in `internal/service/<module>`.
3. Implement repositories in `internal/infra/postgres/<module>`.
4. Add HTTP handlers and route registration in `internal/controller/<module>`.
5. Wire the module in `cmd/api/main.go`.

This structure keeps each module vertically organized while preserving consistent application boundaries.

## Why This Template

This repository is intended to be indexed and reused as a backend template. It favors:

- explicit dependencies
- simple constructor-based wiring
- transaction-safe service methods
- reusable transport and persistence helpers
- a minimal but production-oriented starting point

If you need a Go backend starter with clean architecture principles and PostgreSQL transaction patterns already in place, this template is the base layer.
