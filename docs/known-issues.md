# Known Issues

Defects, drift and risky patterns in the **current** code. Found in the review of 2026-11-05 ([architecture/review.md](architecture/review.md)).

**For agents:** do not copy these patterns into new code. Follow [conventions/](conventions/golang.md) instead. When you fix an issue, delete its entry in the same change.

Severity: **HIGH** = wrong data or broken build/feature · **MEDIUM** = wrong API behaviour or will cause bugs in new code · **LOW** = hygiene.

## Schema drift

<a id="ki-01"></a>

#### KI-01: category table name mismatch `[HIGH]`
- **Where:** `internal/entity/ticket_category.go` `TicketCategory.TableName()` returns `ticket_categories`. Migration `20260319101826_create_table_event_categories.sql` creates `event_category`.
- **Canonical:** `ticket_categories` (decided).
- **Fix:** new migration `ALTER TABLE event_category RENAME TO ticket_categories;`. Do not edit the already-applied migration.

<a id="ki-02"></a>

#### KI-02: order category column mismatch `[HIGH]`
- **Where:** migration `20260319102624_create_table_orders.sql` has `bid_category_id BIGINT`. `entity.Order.TicketCategoryID` maps to `ticket_category_id`.
- **Canonical:** `ticket_category_id` (`bid` was a typo; decided).
- **Fix:** new migration that renames the column, adds FK → `ticket_categories(id)` and an index.

<a id="ki-03"></a>

#### KI-03: `Order.ID` tagged `autoIncrement` `[MEDIUM]`
- **Where:** `internal/entity/order.go:6`.
- **Problem:** the column is a plain `BIGINT PRIMARY KEY` with no sequence. IDs must be Sonyflake ([RULE-PLT-01](invariants/platform.md)). With this tag GORM may omit the ID on insert.
- **Fix:** remove `autoIncrement`; generate the ID in the repository.

<a id="ki-18"></a>

#### KI-18: `tickets.event_id` is nullable `[LOW]`
- **Where:** `20260319102429_create_table_tickets.sql`.
- **Fix:** migration adding `NOT NULL` (a tier always belongs to an event).

## Error handling and API behaviour

<a id="ki-04"></a>

#### KI-04: "not found" returned as `nil, nil` `[MEDIUM]`
- **Where:** `repository.GetByID` (`internal/repository/event.go:89`), `repository.GetListEvent` (`:126`), `service.GetListEvents` (`internal/service/event.go:63`).
- **Problem:** every caller must remember a nil check, or it panics.
- **Rule for new code:** [conventions/golang.md#errors](conventions/golang.md#errors).

<a id="ki-05"></a>

#### KI-05: error → status map uses exact equality `[MEDIUM]`
- **Where:** `response.ResponseError` / `errResponse` (`internal/controller/response/formatter.go:34`).
- **Problem:** `fmt.Errorf("...: %w", ErrEventNotFound)` is not found in the map → 500.
- **Fix:** loop over the map with `errors.Is`. Rule: [conventions/api.md](conventions/api.md#errors-and-status-codes).

<a id="ki-06"></a>

#### KI-06: no input validation `[MEDIUM]`
- **Where / symptoms:**
  - `CreateEventRequestDTO` has no `binding` tags. A missing or empty `name` is saved as `""`. The DB `NOT NULL` does not catch it, because GORM writes the Go zero value.
  - A malformed JSON body → `ShouldBindJSON` error → not mapped → **500** instead of 400 (`internal/controller/event.go:36`).
  - `helper.StringToDateTime` logs a parse error and returns the zero time, which is then saved.
  - `helper.ExpectedNumber[uint64](c.Param("id"))` returns `0` for `"abc"` → 404 instead of 400.
  - No check that `start_time < end_time`.

<a id="ki-12"></a>

#### KI-12: delete of a missing ID returns 200 `[LOW]`
- **Where:** `repository.DeleteByID` (`internal/repository/event.go:134`) ignores `RowsAffected`.

<a id="ki-13"></a>

#### KI-13: empty list returns `"data": null` `[LOW]`
- **Where:** `service.GetListEvents`. New code returns `[]`.

<a id="ki-17"></a>

#### KI-17: response JSON is camelCase, request JSON is snake_case `[LOW]`
- **Where:** `GetListEventsResponseDTO`, `GetEventByIDResponseDTO` (`internal/entity/event.go`) vs `CreateEventRequestDTO`.
- **Decision:** new endpoints use snake_case for both. Existing event responses stay as they are until a versioned change, because changing them breaks clients.

## Build and lint

<a id="ki-08"></a>

#### KI-08: `go vet` fails → `make lint` and `make test` fail `[HIGH]`
- **Where:** `internal/middleware/trace.go:72` `fmt.Errorf(c.Errors.String())` (non-constant format string).
- **Fix:** `errors.New(c.Errors.String())`.

<a id="ki-16"></a>

#### KI-16: golangci-lint findings in `cacher` `[LOW]`
- `cacher/cacher_manager.go:66` `defer gz.Close()` error unchecked (errcheck).
- `cacher/cacher_manager.go:86` if/else on `ttl` → tagged switch (staticcheck QF1003).

## Observability

<a id="ki-09"></a>

#### KI-09: incoming trace context is discarded `[MEDIUM]`
- **Where:** `internal/middleware/trace.go:32`. The `ctx` from `propagator.Extract(...)` is overwritten by `tracing.StartWithName(c.Request.Context(), …)`.
- **Effect:** upstream `traceparent` headers are ignored; every request starts a new trace.
- **Fix:** pass the extracted `ctx` to `StartWithName`.

<a id="ki-22"></a>

#### KI-22: request/response bodies recorded on spans `[HIGH once orders exist]`
- **Where:** `TraceMiddleware` sets `http.request.body` and `http.response.body` attributes.
- **Problem:** order payloads will contain personal data (phone, email, blood type, emergency contact). It would be exported to the tracing backend.
- **Fix before orders ship:** redact or skip bodies for order routes. See [conventions/observability.md](conventions/observability.md).

## Data types and IDs

<a id="ki-07"></a>

#### KI-07: money stored as `float64` `[MEDIUM]`
- **Where:** `Ticket.Price`, `Order.PriceAmount`.
- **Fix:** use `github.com/shopspring/decimal` (already an indirect dependency) or integer minor units. Rule: [conventions/database.md](conventions/database.md#money).

<a id="ki-14"></a>

#### KI-14: event ID generated twice `[LOW]`
- **Where:** `CreateEventRequestDTO.ToEventEntity` (`internal/entity/event.go:54`) generates a Sonyflake ID and ignores the error. `repository.Insert` generates one again only if missing.
- **Rule:** the repository owns ID generation ([conventions/database.md](conventions/database.md#ids)).

## Caching

<a id="ki-15"></a>

#### KI-15: `DeleteByID` deletes from Redis even when caching is off `[LOW]`
- **Where:** `internal/repository/event.go:150`. Harmless, but inconsistent with the reads, which check `config.EnableCaching()`. If Redis is down it logs an error on every delete.

## Structure and testability

<a id="ki-10"></a>

#### KI-10: hidden global state `[LOW]`
- `config.LoadConfig()` runs in `internal/cmd/root.go` `init()` and is fatal without `.env`.
- `utils/singleton` builds Sonyflake in `init()`. `helper.LocationJakarta` is loaded in `init()` with the error ignored (`utils/helper/time.go:14`); if tzdata is missing, it is `nil`.
- Config is read via `os.Getenv` on every call.
- **Effect:** tests that import these packages need the environment set up.

<a id="ki-11"></a>

#### KI-11: no tests `[HIGH]`
- No `_test.go` files exist. The stock logic in particular needs a concurrency test ([RULE-TKT-02](invariants/ticketing.md)).

<a id="ki-19"></a>

#### KI-19: time zone hard-coded to Asia/Jakarta `[LOW]`
- **Where:** `utils/helper/time.go`. Acceptable while every race is in WIB. `DB_TIMEZONE` and this constant must agree.

<a id="ki-20"></a>

#### KI-20: `gin.SetMode` called after `gin.Default()` `[LOW]`
- **Where:** `internal/cmd/server.go`. The engine is created while still in debug mode (debug route logs are printed), and release mode only applies after that.

<a id="ki-21"></a>

#### KI-21: `.env` is committed `[MEDIUM]`
- `.env` is tracked by git. Real credentials must not be committed. Keep only local-dev values, or replace it with `.env.example` and git-ignore `.env`.
