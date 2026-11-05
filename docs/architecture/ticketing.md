# Ticketing Domain — Ticket Tiers, Categories, Orders, Stock

**Status:** Schema and Go entities exist. **No repository, service, controller or route yet.** Everything about behaviour below is **Planned**.
**Key paths:** `internal/entity/ticket.go`, `internal/entity/ticket_category.go`, `internal/entity/order.go`, migrations `20260319101826_*`, `20260319102429_*`, `20260319102624_*`
**Rules:** [../invariants/ticketing.md](../invariants/ticketing.md) — read before writing any code here.
**Vocabulary:** [../domain/glossary.md](../domain/glossary.md)

## Model

```
events 1 ──< tickets (sales tier: price, sale window, QUOTA)
                 1
                 │
                 ^
              orders  >── 1 ticket_categories (race distance, no stock)
```

- **Quota is per ticket tier.** Categories carry no stock. Example: tier *Early Bird* has 500 slots, and a runner using one of them can pick 5K or 10K.
- The price charged is the tier's `price` at order time. It is copied into `orders.price_amount`.

## Data model

Names here are **canonical**: the Go entity names win where SQL and Go disagree. The SQL is fixed by a future migration ([KI-01](../known-issues.md#ki-01), [KI-02](../known-issues.md#ki-02)).

### `tickets` — sales tier
| Column | SQL | Go | Notes |
|--------|-----|----|-------|
| `id` | `BIGINT PK` | `ID uint64` | Sonyflake |
| `event_id` | `BIGINT` FK → `events.id` `ON DELETE CASCADE`, indexed | `EventID` | Nullable in SQL. Should be `NOT NULL`. |
| `name` | `VARCHAR NOT NULL` | `Name` | e.g. "Early Bird" |
| `banner_image` | `VARCHAR` | `BannerImage` | |
| `max_quota` | `INTEGER NOT NULL` | `MaxQuota int` | Total slots in the tier |
| `available_quota` | `INTEGER NOT NULL` | `AvailableQuota int` | Slots left. `0 ≤ available_quota ≤ max_quota` |
| `start_time` / `end_time` | `TIMESTAMP` | `StartTime` / `EndTime` | **Sale window** of the tier (not race time) |
| `price` | `DECIMAL(12,2) NOT NULL` | `Price float64` | Should be a decimal type in Go ([KI-07](../known-issues.md#ki-07)) |
| `created_at` / `updated_at` | `TIMESTAMP` | | |

### `ticket_categories` — race distance
| Column | SQL | Go | Notes |
|--------|-----|----|-------|
| `id` | `BIGINT PK` | `ID uint64` | |
| `name` | `VARCHAR NOT NULL` | `Name` | e.g. "10K" |
| `created_at` / `updated_at` | `TIMESTAMP` | | |

The table is created as `event_category` today ([KI-01](../known-issues.md#ki-01)). It has no `event_id`; whether categories are global or per event is an [open question](../open-questions.md#oq-05).

### `orders` — one runner's registration
| Column | SQL | Go | Notes |
|--------|-----|----|-------|
| `id` | `BIGINT PK` | `ID uint64` | Sonyflake. The Go tag wrongly says `autoIncrement` ([KI-03](../known-issues.md#ki-03)) |
| `ticket_id` | `BIGINT NOT NULL` FK → `tickets.id` `ON DELETE RESTRICT`, indexed | `TicketID` | The tier whose quota this order holds |
| `ticket_category_id` | today `bid_category_id BIGINT` (no FK) | `TicketCategoryID` | Chosen distance ([KI-02](../known-issues.md#ki-02)) |
| `order_id` | `VARCHAR NOT NULL UNIQUE`, indexed | `OrderID` | Public order reference. Format: [open question](../open-questions.md#oq-07) |
| `email`, `phone`, `name`, `gender` | `VARCHAR` | | Runner identity. `email` indexed |
| `bib_number` | `VARCHAR` | `BibNumber` | Assignment rule: [open question](../open-questions.md#oq-04) |
| `blood_type`, `emergency_call`, `jersey_size` | `VARCHAR` | | Race-day data |
| `payment_status` | `VARCHAR NOT NULL DEFAULT 'pending'` | `PaymentStatus` | See lifecycle below |
| `price_amount` | `DECIMAL(12,2) NOT NULL` | `PriceAmount float64` | Snapshot of `tickets.price` at order time |
| `created_at` / `updated_at` | `TIMESTAMP` | | |
| `expired_at` | `TIMESTAMP` | `ExpiredAt *time.Time` | When a `pending` order stops holding its slot |

## Stock lifecycle — Planned (decided by the product owner)

```
                      create order (tier on sale, available_quota > 0)
                      available_quota -= 1   (atomic, same transaction)
                                │
                                ▼
                          ┌──────────┐   payment confirmed    ┌──────┐
                          │ pending  │ ─────────────────────▶ │ paid │  slot stays taken
                          └──────────┘                        └──────┘
                                │
              now ≥ expired_at, still unpaid
              available_quota += 1  (exactly once)
                                ▼
                          ┌──────────┐
                          │ expired  │
                          └──────────┘
```

1. **Reserve on create.** Creating a `pending` order takes one slot from the tier immediately. If `available_quota` is 0, the order is rejected.
2. **Release when it leaves `pending` unpaid.** If the order is not paid by `expired_at`, it becomes expired and the slot returns to `available_quota`, **exactly once**. Expiry is the only such path defined so far. Cancellation and failed payment are [open questions](../open-questions.md#oq-01).
3. **Paid orders keep their slot.** Refunds: [open question](../open-questions.md#oq-01).

Status names `paid` and `expired` are working names. The full set of `payment_status` values is [open question OQ-01](../open-questions.md#oq-01).

Implementation constraints are in [../invariants/ticketing.md](../invariants/ticketing.md). Most important: the decrement must be one conditional SQL statement, so concurrent buyers cannot oversell:

```sql
UPDATE tickets SET available_quota = available_quota - 1, updated_at = now()
WHERE id = $1 AND available_quota > 0
  AND now() BETWEEN start_time AND end_time;   -- 0 rows affected ⇒ sold out / not on sale
```

The order insert and the decrement go in **one DB transaction**, owned by the **order repository** (one method, e.g. `CreateWithReservation`), which also updates `tickets`. Do not do read-check-write in Go, and do not use the Redis cache to decide stock.

## Not built yet (expected next: ticket categories & stock)

"Categories & stock" means **category CRUD + quota on ticket tiers**. Categories never get stock ([RULE-TKT-08](../invariants/ticketing.md)). Whether this feature also includes ordering (create, pay, expire) is [OQ-11](../open-questions.md#oq-11).

- CRUD for ticket tiers per event, and for categories
- Create order (reserve stock), payment confirmation, expiry job (release stock)
- Listing available tiers for an event with remaining quota

Before building these, resolve or explicitly defer the items in [../open-questions.md](../open-questions.md), and write a reqspec in [../reqspecs/](../reqspecs/README.md).
