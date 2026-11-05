# Plan: make go-event-ticketing AI-friendly with a `docs/` single source of truth

## Context

The user wrote this repo before AI coding tools existed and stopped working on it. They now want Claude to do the work here, so the repo has to answer two questions for an AI on its own:

1. **How is the code built?** This is partly answered by CLAUDE.md, which was created this session.
2. **What are the business rules?** These currently exist only in the user's head.

The next feature Claude will build is **ticket categories & stock**, so the docs must make the stock and order rules clear before that work starts.

Reference structure: `/Users/reosahobby/Documents/astro/astro-ads-be/docs`. It has per-domain architecture docs, `invariants/` written as `RULE-XX-NN` with **Invariant / Check for / Violation looks like** and a severity, `conventions/`, and dated `reqspecs/` and `implspecs/` folders. We use a slimmed version of it.

**Scope (user decision):** docs only. Code issues are recorded in the docs, not fixed. Fixing them is separate later work.

## Decisions from the user

- Domain: **running-race registration**.
  - **Event** = one race.
  - **Ticket** = a sales tier for an event (Early Bird, Normal), with its own price, sale window and quota.
  - **Category** = race distance (5K, 10K…). The runner picks one. It has **no stock of its own**.
  - **Order** = one runner's registration (bib, jersey, blood type, emergency contact).
- **Quota is per ticket** (`tickets.max_quota` / `available_quota`).
- **Stock rule:** creating a pending order reserves 1 from `available_quota`. If the order is still unpaid at `expired_at`, it expires and the quota is returned.
- Folder set: slim set (below). Leave out runbooks, monitoring, observability JSON, infra and bug-reports.
- **Conventions tell new code how to be written, not just what exists.** Wherever the current code shows a known issue, the convention gives the better rule. Each convention file has an "As built" note pointing to `known-issues.md`. The better rules:
  - "Not found" is a sentinel error (`ErrXNotFound`), not `nil, nil`.
  - Error mapping uses `errors.Is`.
  - Lists return `[]`.
  - Bad path or body input → 400 through `ErrBadRequest`.
  - Request DTOs use `binding` tags.
  - Money is `decimal`, not `float64`.
- **Go names are canonical:** the table is `ticket_categories` and the column is `orders.ticket_category_id`. `bid_category_id` was a typo, and the `event_category` table name is wrong. A future migration fixes the SQL; that fix is recorded in `known-issues.md`, not made now.

My calls (check these):
- **JSON casing for new endpoints:** snake_case for both requests and responses. It matches the DB and the existing request side. Existing camelCase responses are listed as a known issue, because changing them breaks the API.
- **Stock rules are marked `Planned`** (nothing enforces them yet). Release is worded generally: "quota is returned exactly once when an order leaves `pending` without being paid". Expiry is the only such path defined today; cancellation and failed payment are open questions.

## Architecture review (verdict for the user, also written to `docs/architecture/review.md`)

**Overall, the code is AI-friendly.** The pattern is clear and repeated: controller → service → repository. Interfaces sit in `entity`, there is one composition root (`router.SetupRouter`), the per-method tracing and logging boilerplate is identical everywhere, and packages are small. An AI can copy the `event` domain to add a new one.

What makes it harder for an AI (all recorded in `known-issues.md`):

- **Schema and entity drift.**
  - `TicketCategory.TableName()` is `ticket_categories`, but the migration creates `event_category`.
  - `orders.bid_category_id` (SQL) vs `TicketCategoryID` (Go).
  - `Order.ID` is tagged `autoIncrement`, but the schema expects a Sonyflake ID.
  - Categories have no link to an event.
- **Implicit conventions that are easy to get wrong.**
  - "Not found" is returned as `nil, nil`.
  - The error→HTTP mapping is an exact-match map, so wrapped errors become 500.
  - `ExpectedNumber` turns a non-numeric `:id` into 0, which gives a 404 instead of a 400.
  - Deleting a missing ID returns 200.
  - Lists return `null` instead of `[]`.
- **No input validation.** Request DTOs have no binding tags. `StringToDateTime` logs a bad date and returns the zero time.
- IDs are generated twice (in `ToEventEntity` and again in the repository).
- Money is `float64`. The timezone is hard-coded to `Asia/Jakarta`. Config is global `os.Getenv` plus `init()` singletons.
- `trace.go`: the context extracted from incoming headers is discarded. `fmt.Errorf` is called with a non-constant format string, which fails `go vet` and `make test`.
- The cache delete in `DeleteByID` ignores `ENABLE_CACHING`. `cacher_manager.go` has two lint findings.
- **No tests**, so an AI has no safety net. This is the biggest gap before the stock work, which depends on concurrency.

## Files to create

```
docs/
  README.md                     index: each folder's purpose + "read before…" table
  architecture/
    README.md                   index of architecture docs
    overview.md                 layers, request flow, composition root, infra (Postgres, Redis, OTLP), config, IDs, time zone
    event.md                    event domain: model, endpoints, DTOs, cache keys, behaviour as built
    ticketing.md                ticket tiers, categories, orders: data model from migrations, intended behaviour, stock lifecycle (state diagram), what is not built yet
    review.md                   architecture review above
  invariants/
    README.md                   format + severity legend (copied from reference style)
    ticketing.md                RULE-TKT-xx [Planned]: quota never < 0 and never > max; atomic conditional decrement; reserve on pending create;
                                return quota exactly once when an order leaves pending unpaid (expiry today); sale window enforced;
                                price snapshot into order.price_amount
    platform.md                 RULE-PLT-xx: IDs are Sonyflake (never DB-generated); schema only via migrations; cache is never source of truth
                                and must be invalidated on mutation; sentinel errors must be registered in errResponse
  conventions/
    golang.md                   per-method tracing/logging boilerplate, constructor/interface placement, error flow, nil,nil not-found
    api.md                      route layout /v1/<domain>, response envelope {message,data}, status codes, JSON casing (request snake_case, response camelCase — as built)
    database.md                 migration workflow (create-migration-file, db-migrate), GORM usage, keep entity tags in sync with SQL
    observability.md            tracing.Start/StartWithName, TraceMiddleware attributes, logrus fields
  domain/
    glossary.md                 event, ticket (tier), category (distance), order, bib, quota, reservation, expiry
  reqspecs/README.md            purpose + template (FR-xx, Given/When/Then ACs, open questions)
  implspecs/README.md           purpose + template (dated file, per-step ACs, verification)
  known-issues.md               every finding from the review, with file:line, severity, suggested fix
  open-questions.md             human-brain decisions still missing (below)
docs/superpowers/specs/2026-11-05-ai-friendly-docs-design.md   this design (brainstorming output)
```

Also update `CLAUDE.md`: add a short "Docs" section pointing to `docs/README.md`. Say that `docs/invariants/` must be read before changing tickets or orders, and that docs must be updated together with any code change they describe.

## Open questions to record (not to invent answers for)

These go in `open-questions.md` and are linked from `ticketing.md`:
- Order states other than `pending` (paid / expired / cancelled / refunded?)
- How long a reservation lasts (the `expired_at` window), and what runs the expiry (cron, Redis TTL, lazy check)
- Payment provider and how payment is confirmed
- How bib numbers are assigned
- Whether categories are a global list or per event (no `event_id` today); whether a ticket is limited to certain categories
- Max tickets per email or order
- Whether the `order_id` string format is defined anywhere

## Writing rules

- Every fact about the code must cite a `path:line` or a symbol and be checked against the source. Anything not built yet is clearly labelled **Planned** or **Open question**.
- Write for an AI reader: short sections, tables, rule IDs. No filler.
- English throughout (code comments in Indonesian are left alone).

## Verification

- Re-read each doc against the code it cites: table and column names come from the migrations, and endpoints from `internal/router/router.go`.
- Check that relative links resolve (`grep -o '](.*\.md)'` and test each path).
- Run a fresh subagent test: give it only `docs/` + `CLAUDE.md` and ask it to outline how it would build "ticket categories & stock". It should cite RULE-TKT rules and the atomic decrement without being told.
- Nothing is committed unless the user asks.
