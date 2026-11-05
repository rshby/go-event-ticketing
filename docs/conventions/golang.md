# Go Conventions

How **new** Go code in this repo is written. Where existing code differs, the difference is listed under *As built* and in [../known-issues.md](../known-issues.md). Do not copy the old pattern.

## Adding a domain — file layout

For a domain `ticket`:

| File | Contains |
|------|----------|
| `internal/entity/ticket.go` | GORM model + `TableName()`, `TicketRepository` and `TicketService` interfaces, request/response DTOs, DTO → entity mappers |
| `internal/repository/ticket.go` | `type ticketRepository struct{ db *gorm.DB; cache cacher.CacheManager }`, `NewTicketRepository(...) entity.TicketRepository` |
| `internal/service/ticket.go` | `type ticketService struct{...}`, `NewTicketService(...) entity.TicketService` |
| `internal/controller/ticket.go` | `type TicketController struct{...}`, `NewTicketController(...) *TicketController` |
| `internal/router/router.go` | wire repo → service → controller in `SetupRouter`, add the route group |
| `internal/errors/errors.go` + `controller/response/formatter.go` | new sentinels + status mapping |

Rules:
- Constructors return the **interface** from `entity` (repository, service), except controllers, which return the concrete struct.
- Implementations are unexported.
- Every method takes `ctx context.Context` first (controllers take `*gin.Context`).
- Each layer depends only on the layer below it, through `entity` interfaces.

## Method boilerplate

Every controller, service, repository and cacher method starts like this:

```go
func (s *ticketService) ReserveSlot(ctx context.Context, ticketID uint64) error {
	ctx, span := tracing.Start(ctx) // span name = this function's full name
	defer span.End()

	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context":   helper.DumpIncomingContext(ctx),
		"ticket_id": ticketID,
	})
	...
}
```

- Always reassign `ctx` from `tracing.Start` and pass that `ctx` down.
- Log an error **once, where it is handled or first seen**. Then return it unchanged, or wrapped with `%w`.
- Do not put personal data in log fields ([RULE-TKT-09](../invariants/ticketing.md)).

## Errors

**Rule for new code:**
- "Not found" is a **sentinel error**, not `nil, nil`. The repository returns `nil, errors.ErrTicketNotFound` when `errors.Is(err, gorm.ErrRecordNotFound)`. Callers never nil-check a result whose error was nil.
- Define sentinels in `internal/errors/errors.go` with `errors.New`. Name them `Err<Thing><Condition>` (`ErrTicketSoldOut`, `ErrTicketNotOnSale`).
- Wrap with `fmt.Errorf("reserve slot: %w", err)` when adding context. The mapping uses `errors.Is` (see [api.md](api.md#errors-and-status-codes)).
- Never `panic` in request paths.

*As built:* the event repository returns `nil, nil` for not found ([KI-04](../known-issues.md#ki-04)), and the status map uses exact equality ([KI-05](../known-issues.md#ki-05)). Until KI-05 is fixed, **return sentinels unwrapped** from the service to the controller.

## Lists

Return an empty slice, never nil: `make([]T, 0, len(rows))`. JSON then shows `[]`.
*As built:* `GetListEvents` returns nil ([KI-13](../known-issues.md#ki-13)).

## Types

- Money: decimal, never `float64` ([database.md](database.md#money)).
- IDs: `uint64` internally, string in JSON.
- Time: `time.Time` in Jakarta time from `helper.GetNowTime()`. Parse input with a function that **returns an error**. Never use a zero time silently ([KI-06](../known-issues.md#ki-06)).

## Transactions

Multi-step writes (e.g. order insert + quota decrement) use `db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {...})` inside the repository. The service does not handle `*gorm.DB`. When an operation spans several repository methods, add a repository method that does the whole unit of work.

## Tooling

`make fmt` (gofmt -s + goimports), `make lint` (go vet + golangci-lint), `make test` (race + cover), `make build`. Imports are grouped as stdlib, then everything else.
