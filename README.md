# go-event-ticketing

Backend service for **running-race registration**: organisers create races (events) with ticket tiers (e.g. Early Bird, Normal), and runners register for a race distance (5K, 10K, …).

> **Status:** event management is built. Ticket tiers, categories, orders and stock have a database schema but no API yet. See [docs/architecture/ticketing.md](docs/architecture/ticketing.md).

## Tech stack

| Concern | Choice |
|---------|--------|
| Language | Go 1.25 |
| HTTP | [Gin](https://github.com/gin-gonic/gin) |
| Database | PostgreSQL via [GORM](https://gorm.io) |
| Migrations | [sql-migrate](https://github.com/rubenv/sql-migrate) (plain SQL files) |
| Cache | Redis ([go-redis](https://github.com/redis/go-redis)), gzip-compressed JSON |
| Tracing | OpenTelemetry → OTLP gRPC |
| Logging | [logrus](https://github.com/sirupsen/logrus) |
| CLI | [Cobra](https://github.com/spf13/cobra) |
| IDs | [Sonyflake](https://github.com/sony/sonyflake) |

## Getting started

### Prerequisites

- Go 1.25+
- PostgreSQL
- Redis
- An OTLP-compatible trace collector (e.g. Jaeger or the OpenTelemetry Collector) on the configured gRPC endpoint
- Optional, for `make fmt` / `make lint`: `goimports`, `golangci-lint`

### Configure

All settings come from a `.env` file in the repo root, which is loaded on startup and is required. Edit the values for your environment:

| Group | Variables |
|-------|-----------|
| App | `PORT`, `MODE` |
| PostgreSQL | `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_TIMEZONE`, `DB_CONNECTION_TIMEOUT`, `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME`, `DB_CONN_MAX_IDLETIME` |
| Redis | `ENABLE_CACHING`, `REDIS_HOST`, `REDIS_PORT`, `REDIS_DB_NUMBER`, `REDIS_MAX_CONN_SIZE`, `REDIS_IDLE_CONN_SIZE`, `REDIS_CONN_LIFETIME` |
| Tracing | `OTLP_SERVICE_NAME`, `OTLP_ENDPOINT`, `OTLP_PORT` |
| HTTP server | `HTTP_SERVER_READ_HEADER_TIMEOUT`, `HTTP_SERVER_READ_TIMEOUT`, `HTTP_SERVER_WRITE_TIMEOUT`, `HTTP_SERVER_IDLE_TIMEOUT` |

Durations use Go syntax (`5s`, `10m`). Set `DB_TIMEZONE=Asia/Jakarta`, which is the business time zone used by the app. Optional `MACHINE_ID` sets the Sonyflake machine ID (otherwise it is derived from the private IP).

### Run

Run every command from the repo root (that's where `.env` is read from).

```bash
go run . db-migrate      # create / update the database schema
go run . server          # start the HTTP API on $PORT
```

## Commands

| Command | Description |
|---------|-------------|
| `go run . server` | Start the HTTP server (graceful shutdown on SIGINT/SIGTERM) |
| `go run . db-migrate` | Apply pending migrations from `internal/database/migration` |
| `go run . create-migration-file <name>` | Create a new timestamped migration file |
| `make fmt` | Format code (`gofmt -s`, `goimports`) |
| `make lint` | `go vet` + `golangci-lint` |
| `make test` | Run tests with the race detector and coverage |
| `make build` | Build the binary to `bin/go-event-ticketing` |
| `make all` | fmt, lint, test, build |

## API

Base path: `/v1`. Every response uses the same envelope:

```json
{ "message": "success retrieves event by id", "data": { } }
```

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/v1/event` | Create an event |
| `GET` | `/v1/event` | List events (newest first) |
| `GET` | `/v1/event/:id` | Get an event by ID |
| `DELETE` | `/v1/event/:id` | Delete an event |

Create an event:

```bash
curl -X POST http://localhost:$PORT/v1/event \
  -H 'Content-Type: application/json' \
  -d '{
        "name": "Jakarta Marathon 2026",
        "banner_image": "https://example.com/banner.jpg",
        "location": "Jakarta",
        "start_time": "2026-12-01 05:00:00",
        "end_time": "2026-12-01 12:00:00"
      }'
```

Date-times use the format `YYYY-MM-DD hh:mm:ss` in Asia/Jakarta time. IDs are returned as strings.

## Project structure

```
main.go                 entry point → Cobra CLI
internal/cmd/           CLI subcommands (server, db-migrate, create-migration-file)
internal/router/        route registration and dependency wiring
internal/controller/    HTTP handlers and response formatting
internal/service/       business logic
internal/repository/    database access and caching
internal/entity/        models, interfaces and DTOs per domain
internal/database/      Postgres/Redis connections and SQL migrations
internal/middleware/    tracing middleware
internal/errors/        sentinel errors
cacher/                 Redis cache manager and keys
tracing/                OpenTelemetry setup and span helpers
config/                 environment configuration
utils/                  helpers and the Sonyflake ID generator
docs/                   project documentation (single source of truth)
```

## Documentation

The [docs/](docs/README.md) folder is the single source of truth for architecture, business rules and conventions:

- [Glossary](docs/domain/glossary.md): what event, ticket tier, category and order mean
- [Architecture overview](docs/architecture/overview.md)
- [Invariants](docs/invariants/README.md): rules that must never be broken
- [Conventions](docs/conventions/golang.md): how new code is written
- [Known issues](docs/known-issues.md) and [open questions](docs/open-questions.md)

AI coding agents should start with [CLAUDE.md](CLAUDE.md).
