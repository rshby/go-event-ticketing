# API Conventions

## Routes

- Every route is under `/v1`, grouped by domain in `router.SetupRouter`.
- The group name is the **singular** resource: `/v1/event`, `/v1/ticket`, `/v1/order`.
- Standard shape:

| Action | Method + path |
|--------|---------------|
| Create | `POST /v1/<resource>` |
| List | `GET /v1/<resource>` |
| Get | `GET /v1/<resource>/:id` |
| Update | `PUT /v1/<resource>/:id` |
| Delete | `DELETE /v1/<resource>/:id` |
| Child list | `GET /v1/event/:id/ticket` (e.g. tiers of one event) |
| State change | `POST /v1/<resource>/:id/<verb>` (e.g. `/v1/order/:id/pay`) |

## Response envelope

Every response, success or error, has this shape (`response.ResponseDTO`, `internal/controller/response/formatter.go`):

```json
{ "message": "success create ticket", "data": { ... } }
```

- Success: `response.NewResponse().WithMessage("...").WithData(&x).ToResponseAPI(c, status)`.
- Error: `response.ResponseError(c, err)`. `data` is `null`.
- Message wording: `"success <verb> <resource>"`.

## JSON

- **New endpoints: snake_case for request and response fields** (`available_quota`, `start_time`).
  *As built:* event responses use camelCase ([KI-17](../known-issues.md#ki-17)). Leave them as they are; don't extend them.
- IDs are **strings** in JSON (`"id": "123456789"`), because Sonyflake IDs exceed JavaScript's safe integer range.
- Date-times are `"YYYY-MM-DD hh:mm:ss"` in Asia/Jakarta (`time.DateTime`).
- Money: a string with 2 decimals (`"150000.00"`), which avoids float rounding in clients.
- Lists: `"data": []` when empty, never `null`.

## Validation

- Request DTOs carry Gin `binding` tags (`binding:"required"`, `gt=0`, …). Bind with `c.ShouldBindJSON`. A bind or validation failure → `ErrBadRequest` → 400.
- Path IDs: parse with `strconv.ParseUint`. Failure → 400. Don't use `helper.ExpectedNumber` for input, because it hides errors ([KI-06](../known-issues.md#ki-06)).
- Cross-field rules (start before end, quota ≥ 1) go in the service and return a sentinel → 400.
- Never trust client-sent prices, quotas or statuses on order creation ([RULE-TKT-07](../invariants/ticketing.md)).

## Errors and status codes

| Situation | Status | Error |
|-----------|--------|-------|
| Malformed body, validation failure, bad path param | 400 | `ErrBadRequest` or a specific sentinel |
| Resource does not exist (GET, PUT, DELETE) | 404 | `Err<Thing>NotFound` |
| Business conflict (sold out, not on sale, already paid, expired) | 409 | `ErrTicketSoldOut`, `ErrTicketNotOnSale`, … |
| Anything else | 500 | message `"internal error"`. The real error is logged, never returned |

Mapping: add each sentinel to `errResponse` in `formatter.go` ([RULE-PLT-05](../invariants/platform.md)). Matching with `errors.Is` needs the [KI-05](../known-issues.md#ki-05) fix. Fix KI-05 first, in the same change as the first new domain. If it isn't fixed yet, return sentinels unwrapped.

Success codes: 201 for create, 200 otherwise.
