# Habit Tracker

Backend service for habit tracking. The application provides HTTP APIs for users, habits, streaks and achievements, persists data in PostgreSQL, and uses Kafka with an outbox publisher for asynchronous metric and achievement processing.

## Stack

- Go 1.25
- Gin HTTP router
- PostgreSQL
- Kafka
- Docker Compose
- golangci-lint
- GitHub Actions CI

## Project Structure

```text
cmd/server              application entrypoint
config                  configuration and logger setup
internal/adapter/http   HTTP handlers, routers, DTOs and middleware
internal/adapter/kafka  Kafka producer and consumer adapters
internal/adapter/postgres PostgreSQL repositories and tx manager
internal/domain         domain models and errors
internal/usecase        application services
internal/worker/outbox  outbox publisher worker
migrations              database migrations
tests/e2e               end-to-end tests
```

## Requirements

- Go 1.25+
- Docker and Docker Compose
- Make

On Windows, `make test` may fail in `tests/e2e` if Docker/testcontainers cannot create a supported provider. The internal packages can still be checked with:

```bash
go test ./cmd/... ./config/... ./internal/...
```

## Configuration

Copy `.env.example` to `.env` and adjust values if needed:

```bash
cp .env.example .env
```

Important variables:

- `SERVER_HOST`, `SERVER_PORT`
- `POSTGRES_HOST`, `POSTGRES_PORT`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_SSL_MODE`
- `JWT_SECRET`
- `KAFKA_BROKERS`
- `KAFKA_OUTBOX_TOPIC`
- `OUTBOX_PUBLISH_INTERVAL`

Docker Compose overrides some values for container networking, for example PostgreSQL host is `habit-tracker-db` and Kafka broker is `kafka:29092` inside the compose network.

## Running Locally

Start all services:

```bash
make docker-up
```

The API will be available at:

```text
http://localhost:8080
```

Kafka UI is available at:

```text
http://localhost:8085
```

Useful commands:

```bash
make logs            # tail app logs
make docker-logs     # tail all compose logs
make db-shell        # open psql in the database container
make docker-down     # stop compose services
```

Run the app without Docker:

```bash
make run
```

This requires PostgreSQL and Kafka to be reachable from values in `.env`.

## Makefile Commands

```bash
make help             # show available targets
make run              # run app locally
make build            # build server binary into bin/server
make test             # run all Go tests
make fmt              # format Go files
make tidy             # run go mod tidy
make docker-up        # start docker compose services
make docker-rebuild   # rebuild and start app service
make docker-down      # stop docker compose services
make logs             # tail app logs
make docker-logs      # tail all compose logs
make db-shell         # open psql in Postgres container
make migrate-status   # show migration version and outbox table status
make clean            # remove local build artifacts
```

## HTTP API

Health/info:

```http
GET /_info
```

Auth:

```http
POST /api/register
POST /api/login
POST /api/refresh
```

Habits:

```http
GET    /api/habit
GET    /api/habit/:id
POST   /api/habit
PUT    /api/habit/:id
DELETE /api/habit/:id
POST   /api/habit/:id/add
GET    /api/habit/my
POST   /api/habit/confirm/:id
GET    /api/habit/heatmap
```

Tags:

```http
GET    /api/habit/tag/all
GET    /api/habit/tag/:id
POST   /api/habit/tag
PUT    /api/habit/tag/:id
DELETE /api/habit/tag/:id
```

Achievements:

```http
GET /achievements?page=1&page_size=20
GET /api/achievements?page=1&page_size=20
```

The achievements endpoint returns all achievements and marks user progress:

```json
{
  "achievements": [
    {
      "id": "uuid",
      "code": "first_habit",
      "title": "First Habit",
      "description": "Create your first habit",
      "enabled": true,
      "unlocked": true,
      "unlocked_at": "2026-06-05T12:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 1
  }
}
```

Most user-facing endpoints require `Authorization: Bearer <token>`.

## Events and Achievements

The service uses the outbox pattern:

1. Usecases write domain events into the `outbox_event` table in the same transaction as business changes.
2. `internal/worker/outbox` publishes pending events to Kafka.
3. Metric consumers read habit/streak events and update user metrics.
4. Achievement consumers read metric update events and unlock achievements when conditions are satisfied.

The default compose Kafka topic is:

```text
habit.events.v1
```

## Migrations

Migrations are applied on application startup from the `migrations` directory.

To inspect migration state in Docker:

```bash
make migrate-status
```

## CI

GitHub Actions is configured in `.github/workflows/main.yml`.

The CI pipeline runs on pushes and pull requests to `main` and `dev`:

1. Checkout repository.
2. Set up Go from `go.mod`.
3. Download dependencies with `go mod download`.
4. Run `golangci-lint` for `./...`.
5. Run tests with `go test -v ./...`.
6. Build the server with `go build -v ./cmd/server`.
7. Check Docker image build with `docker/build-push-action`.

Before opening a PR, run:

```bash
make fmt
make test
make build
```
