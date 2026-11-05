# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

The binary is a Cobra CLI (`internal/cmd`); every subcommand loads `.env` from the current working directory and exits fatally if it is missing, so run from the repo root.

```bash
go run . server                               # start HTTP API (needs Postgres, Redis, OTLP collector)
go run . db-migrate                           # apply pending migrations (sql-migrate, table "migrations")
go run . create-migration-file <name>         # scaffold internal/database/migration/<timestamp>_<name>.sql
make fmt                                      # gofmt -s + goimports
make lint                                     # go vet + golangci-lint (default linters, no config file)
make test                                     # go test -race -cover ./... (no tests exist yet)
make build                                    # binary at bin/go-event-ticketing
make all                                      # fmt lint test build
go test ./internal/service -run TestName      # single test
```

There is no Dockerfile or golangci-lint config.

## Docs (single source of truth)

`docs/` holds the business rules and decisions that the code can't show. Start at `docs/README.md`.
- Read `docs/domain/glossary.md` first. This service sells **running-race registrations**: a "ticket" is a sales tier (Early Bird) with quota, and a "category" is a race distance with no stock.
- Before touching tickets, orders or stock, read `docs/invariants/ticketing.md`. Those rules are binding.
- New code follows `docs/conventions/`. Where existing code differs, it is listed in `docs/known-issues.md`. Don't copy those patterns.
- Never invent answers to `docs/open-questions.md`. Ask instead.
- Update the affected docs in the same change as the code.

## Architecture

Go 1.25 HTTP service: Gin + GORM (PostgreSQL) + go-redis + OpenTelemetry (OTLP gRPC), logging via logrus.

**Request flow:** `router.SetupRouter` (`internal/router/router.go`) does all dependency wiring by hand — repository → service → controller — then mounts routes under `/v1`. Add new domains by wiring them there.

**Layering convention (follow `event` as the template):**
- `internal/entity/<domain>.go` holds the GORM model, the `Repository` and `Service` interfaces, and request/response DTOs together. Services/repositories return interfaces from `entity`; controllers are concrete structs.
- `internal/repository` — GORM access plus Redis read-through caching. Caching on reads is gated by `config.EnableCaching()`; cache keys are built in `cacher/cache_key.go` (`ticketing:<MODE>:<key>`). Invalidate the cache on mutations.
- `internal/service` — maps entities to response DTOs (IDs rendered as strings, times via `helper.DateTimeToString`). *As built,* "not found" is signalled by the repository returning `nil, nil`, and the service converts it to a sentinel error. New code returns the sentinel straight from the repository (`docs/conventions/golang.md`).
- `internal/controller` — binds input, calls service, responds via `response.NewResponse().WithMessage(...).WithData(...).ToResponseAPI(c, status)` or `response.ResponseError(c, err)`.

**Errors → HTTP status:** sentinel errors live in `internal/errors/errors.go`; `internal/controller/response/formatter.go` maps them to status codes via the `errResponse` map (exact equality lookup, not `errors.Is`). Any unmapped error becomes a 500 "internal error". New sentinel errors need an entry in that map.

**Per-function boilerplate:** every controller/service/repository/cacher method starts with
```go
ctx, span := tracing.Start(ctx)   // span name auto-derived from caller func name
defer span.End()
logger := logrus.WithContext(ctx).WithFields(logrus.Fields{"context": helper.DumpIncomingContext(ctx), ...})
```
`tracing.Start` accepts a `*gin.Context` and rewrites its request context. `middleware.TraceMiddleware` creates the root span and records request/response bodies as attributes. GORM is traced via the otel plugin.

**IDs:** primary keys are `uint64` Sonyflake IDs (`helper.GenerateSonyFlakeID`, singleton in `utils/singleton`, machine ID from `MACHINE_ID` env or lower 16 bits of private IP), not DB-generated.

**Cache:** `cacher.CacheManager` gzip-compresses JSON before storing. `Store` with `cacher.NewItem(v)` uses default TTL (`config.DefaultRedisTTL`, 30m); `WithNoTTL()` / `WithCustomTTL(d)` override.

**Schema:** managed only by SQL files in `internal/database/migration` (sql-migrate `-- +migrate Up/Down` format) — no GORM AutoMigrate. Tables: `events`, `tickets`, `orders`, and categories, created as `event_category` but canonically `ticket_categories` (see `docs/known-issues.md` KI-01).

**Config:** `config/config.go` exposes one getter per env var (see `.env` for the full list) with fallbacks in `config/default.go`. Add new settings as getters there rather than reading `os.Getenv` elsewhere.

Code comments and log messages are a mix of English and Indonesian.
