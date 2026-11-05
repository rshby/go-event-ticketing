# Architecture Overview

A single Go HTTP service that registers runners for running races. It is one binary with several Cobra subcommands.

| Concern | Choice |
|---------|--------|
| Language | Go 1.25 (`go.mod`) |
| HTTP | Gin (`github.com/gin-gonic/gin`) |
| Database | PostgreSQL via GORM (`gorm.io/gorm`, `gorm.io/driver/postgres`) |
| Migrations | `github.com/rubenv/sql-migrate`, SQL files only |
| Cache | Redis (`github.com/redis/go-redis/v9`), values gzip-compressed JSON |
| Tracing | OpenTelemetry, OTLP gRPC exporter; GORM traced via `gorm.io/plugin/opentelemetry` |
| Logging | logrus, stdout |
| CLI | Cobra (`github.com/spf13/cobra`) |
| IDs | Sonyflake (`github.com/sony/sonyflake/v2`) |

## Entry points

`main.go` → `cmd.Execute()` (`internal/cmd/root.go`). The root command's `init()` sets up the logger and loads `.env` (`config.LoadConfig`, `config/config.go`). A missing `.env` is fatal, so always run from the repo root.

| Subcommand | File | Does |
|------------|------|------|
| `server` | `internal/cmd/server.go` | Connects OTLP, Redis and Postgres, builds the Gin app, serves HTTP, shuts down gracefully on SIGINT/SIGTERM (3 s timeout) |
| `db-migrate` | `internal/cmd/db_migrate.go` | Applies pending `Up` migrations from `internal/database/migration`, tracking them in table `migrations` |
| `create-migration-file <name>` | `internal/cmd/create_migration_file.go` | Writes `internal/database/migration/<YYYYMMDDhhmmss>_<name>.sql` with the sql-migrate template |

## Layers and request flow

```
HTTP request
  → middleware.TraceMiddleware          root span, records request/response body   internal/middleware/trace.go
  → controller.<Domain>Controller       bind input, call service, write response   internal/controller/
  → service (entity.<Domain>Service)    business logic, entity → response DTO      internal/service/
  → repository (entity.<Domain>Repository)  GORM + Redis cache                     internal/repository/
  → PostgreSQL / Redis
```

| Package | Role |
|---------|------|
| `internal/entity` | GORM models, the `Repository` and `Service` **interfaces**, and request/response DTOs, one file per domain. The centre of the dependency graph: every other layer imports it. |
| `internal/controller` | Gin handlers (concrete structs). `controller/response` holds the response envelope and the error → HTTP status mapping. |
| `internal/service` | Unexported structs implementing `entity.<Domain>Service`. |
| `internal/repository` | Unexported structs implementing `entity.<Domain>Repository`. |
| `internal/router` | **Composition root.** `SetupRouter` (`internal/router/router.go:15`) builds repository → service → controller by hand and registers routes. There is no DI container. |
| `internal/errors` | Sentinel errors and `InternalError` (code + message). |
| `internal/database` | `ConnectPostgreSql`, `ConnectRedis`, and the `migration/` SQL files. |
| `cacher` | `CacheManager` (Redis wrapper with gzip+JSON), `Item` (value + TTL), cache key builders. |
| `tracing` | OTLP setup (`common.go`) and span helpers `Start` / `StartWithName` (`span.go`). |
| `config` | One getter per environment variable plus defaults (`default.go`). |
| `utils/helper` | JSON dump helpers, `ExpectedNumber` / `ExpectedString` converters, Sonyflake ID generation, Jakarta time helpers. |
| `utils/singleton` | Package-level Sonyflake instance created in `init()`. |

## Infrastructure

### PostgreSQL
- Connection: `database.ConnectPostgreSql` (`internal/database/postgresql.go`). `PrepareStmt: true`, `SkipDefaultTransaction: true`, `TranslateError: true`, GORM SQL logging at Info level.
- Pool settings come from `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME`, `DB_CONN_MAX_IDLETIME`.
- Schema: created **only** by SQL migrations. No `AutoMigrate`. See [../conventions/database.md](../conventions/database.md).

| Table | Migration | Used by code |
|-------|-----------|--------------|
| `events` | `20260318195001_create_table_events.sql` | Yes (event domain) |
| `event_category` | `20260319101826_create_table_event_categories.sql` | No — the Go entity expects `ticket_categories` (see [known-issues](../known-issues.md#ki-01)) |
| `tickets` | `20260319102429_create_table_tickets.sql` | Entity only, no repository yet |
| `orders` | `20260319102624_create_table_orders.sql` | Entity only, no repository yet |

### Redis
- Connection: `database.ConnectRedis` (`internal/database/redis.go`). DB number defaults to 3.
- Used as a read-through cache in repositories. Turned on by `ENABLE_CACHING=true`. **Never a source of truth** ([RULE-PLT-03](../invariants/platform.md)).
- Keys: `ticketing:<MODE>:<domain key>` built by `cacher.createCacheKey` (`cacher/cache_key.go`).
- Values: JSON then gzip (`cacheManager.Compress`). TTL: `cacher.NewItem(v)` → default 30 min (`config.DefaultRedisTTL`); `cacher.WithNoTTL()` → no expiry; `cacher.WithCustomTTL(d)`.

### OpenTelemetry
- `tracing.ConnectOTLPTrace` (`tracing/common.go`): gRPC exporter to `OTLP_ENDPOINT:OTLP_PORT` without TLS, `AlwaysSample`, W3C TraceContext + Baggage propagators.
- The `server` command fails at startup if the exporter cannot be created.

## Configuration

All config is environment variables loaded from `.env` by `godotenv`. Read them only through getters in `config/config.go`. Never call `os.Getenv` elsewhere. Defaults live in `config/default.go` and inside the getters.

| Group | Variables |
|-------|-----------|
| App | `PORT`, `MODE` (used in cache key prefix) |
| Postgres | `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_TIMEZONE`, `DB_CONNECTION_TIMEOUT`, `DB_MAX_OPEN_CONNS`, `DB_MAX_IDLE_CONNS`, `DB_CONN_MAX_LIFETIME`, `DB_CONN_MAX_IDLETIME` |
| Redis | `ENABLE_CACHING`, `REDIS_HOST`, `REDIS_PORT`, `REDIS_DB_NUMBER`, `REDIS_MAX_CONN_SIZE`, `REDIS_IDLE_CONN_SIZE`, `REDIS_CONN_LIFETIME` |
| OTLP | `OTLP_SERVICE_NAME`, `OTLP_ENDPOINT`, `OTLP_PORT` |
| HTTP server | `HTTP_SERVER_READ_HEADER_TIMEOUT`, `HTTP_SERVER_READ_TIMEOUT`, `HTTP_SERVER_WRITE_TIMEOUT`, `HTTP_SERVER_IDLE_TIMEOUT` |
| IDs | `MACHINE_ID` (optional, see below) |

## IDs

- Every primary key is a `uint64` **Sonyflake ID generated in the application**, never by the database ([RULE-PLT-01](../invariants/platform.md)).
- Generator: `helper.GenerateSonyFlakeID()` → `singleton.SF` (`utils/singleton/singleton.go`). Epoch is 2026-01-01 UTC. The machine ID comes from `MACHINE_ID`, otherwise from the lower 16 bits of the first private IPv4 address.
- IDs are serialised as **strings** in JSON responses (`helper.ExpectedString`). This keeps JavaScript clients from losing precision on 64-bit integers.

## Time

- Business time zone is `Asia/Jakarta`, hard-coded in `helper.LocationJakarta` (`utils/helper/time.go`).
- API date-time format is Go's `time.DateTime` layout (`2006-01-02 15:04:05`), parsed in Jakarta time. There is no time-zone offset in the string.
- DB columns are `TIMESTAMP` (without time zone). The DSN sets `TimeZone=<DB_TIMEZONE>`.
