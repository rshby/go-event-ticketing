# Event Domain

**Status:** Built
**Key paths:** `internal/entity/event.go`, `internal/controller/event.go`, `internal/service/event.go`, `internal/repository/event.go`, `cacher/cache_key.go`
**Business meaning:** one running race. See [../domain/glossary.md](../domain/glossary.md).

This is the **reference implementation**. New domains copy its structure, but follow [../conventions/](../conventions/golang.md) wherever this code has a [known issue](../known-issues.md).

## Data model — table `events`

| Column | SQL type | Go field | Notes |
|--------|----------|----------|-------|
| `id` | `BIGINT PK` | `ID uint64` | Sonyflake |
| `name` | `VARCHAR NOT NULL` | `Name string` | |
| `banner_image` | `VARCHAR` | `BannerImage string` | URL/path, not validated |
| `location` | `VARCHAR` | `Location string` | |
| `start_time` | `TIMESTAMP` | `StartTime time.Time` | |
| `end_time` | `TIMESTAMP` | `EndTime time.Time` | |
| `created_at` | `TIMESTAMP DEFAULT now` | `CreatedAt` | `autoCreateTime` |
| `updated_at` | `TIMESTAMP DEFAULT now` | `UpdatedAt` | `autoUpdateTime` |
| `created_by` | `VARCHAR` | `CreatedBy string` | Always `"admin"` today; no auth exists |

Deleting an event **cascades** to its `tickets` (FK `fk_tickets_event_id`). Orders reference tickets with `ON DELETE RESTRICT`, so deleting an event that has orders **fails at the DB**. Today that surfaces as a 500.

## Endpoints (`internal/router/router.go:26-29`)

| Method | Path | Handler | Success | Body / response |
|--------|------|---------|---------|-----------------|
| `POST` | `/v1/event` | `CreateEvent` | 201 | Request `CreateEventRequestDTO` (snake_case). No data in the response. |
| `GET` | `/v1/event` | `GetListEvents` | 200 | `[]GetListEventsResponseDTO` (camelCase), ordered `id DESC`, no pagination |
| `GET` | `/v1/event/:id` | `GetEventByID` | 200 | `GetEventByIDResponseDTO` (camelCase) |
| `DELETE` | `/v1/event/:id` | `DeleteEventByID` | 200 | No data |

Create request:
```json
{ "name": "Jakarta Marathon", "banner_image": "https://…", "location": "Jakarta",
  "start_time": "2026-12-01 05:00:00", "end_time": "2026-12-01 12:00:00" }
```

Response envelope (all endpoints): `{"message": string, "data": any}`. See [../conventions/api.md](../conventions/api.md).

## Behaviour as built

| Case | Result | Where |
|------|--------|-------|
| Malformed JSON body | 500 "internal error" (the bind error is not mapped) | `controller.CreateEvent`, `response.ResponseError` |
| Bad `start_time` format | Saved as zero time; error only logged | `helper.StringToDateTime` |
| Non-numeric `:id` | Converted to `0` → 404 "event is not found" | `helper.ExpectedNumber` |
| Unknown `:id` on GET | Repo returns `nil, nil` → service returns `ErrEventNotFound` → 404 | `repository.GetByID`, `service.GetEventByID` |
| Unknown `:id` on DELETE | 200 "success delete event" (zero rows is not checked) | `repository.DeleteByID` |
| Empty list | 200 with `"data": null` | `service.GetListEvents` |

These are recorded as [known issues](../known-issues.md). New code must not copy them.

## Caching

| Key | Value | TTL | Written | Invalidated |
|-----|-------|-----|---------|-------------|
| `ticketing:<MODE>:event:id:<id>` (`cacher.GetEventByIDCacheKey`) | `entity.Event` | 30 min default | `GetByID` after a DB hit, only when `ENABLE_CACHING=true` | `DeleteByID` (always, even when caching is off) |

The list endpoint is not cached. Create and update do not touch the cache (there is no update endpoint yet).

## Not built

- Update event endpoint
- Pagination and filtering on the list
- Authentication; `created_by` is hard-coded
