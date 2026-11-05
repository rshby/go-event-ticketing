# Architecture Review — AI-Friendliness

**Date:** 2026-11-05
**Scope:** whole repository at commit `467f537` (develop)
**Question:** is the architecture easy for an AI coding agent to understand and extend safely?

## Verdict

**Yes, the structure is AI-friendly.** The main risks are not the layout. They are **implicit rules** and **schema drift**, which an agent will copy or trip over. Those are now written down in [../known-issues.md](../known-issues.md), [../invariants/](../invariants/README.md) and [../conventions/](../conventions/golang.md).

## What works well for an agent

| Strength | Why it helps |
|----------|--------------|
| One pattern, repeated: controller → service → repository | An agent can add a domain by copying `event` |
| Interfaces and DTOs live together in `internal/entity/<domain>.go` | One file shows a domain's whole contract |
| Single composition root, `router.SetupRouter` | All wiring is in one place; no DI magic to trace |
| Identical per-method boilerplate (span + logger) | Easy to reproduce; consistent traces |
| Small packages, short functions, no generics-heavy abstractions | Low reading cost |
| SQL-only migrations with a generator command | Schema history is explicit and reviewable |
| Config only through getters in `config/` | One place to find every setting |

## What makes it harder (and where it is now documented)

| Risk | Effect on an agent | Documented in |
|------|-------------------|---------------|
| Go entities and SQL disagree on table and column names | Agent writes queries against tables that don't exist | [KI-01](../known-issues.md#ki-01), [KI-02](../known-issues.md#ki-02), [KI-03](../known-issues.md#ki-03) |
| "Not found" returned as `nil, nil` | Agent forgets the nil check → nil pointer panic | [conventions/golang.md](../conventions/golang.md#errors), [KI-04](../known-issues.md#ki-04) |
| Error → HTTP map uses exact equality | Wrapping an error with `%w` silently turns a 404 into a 500 | [conventions/api.md](../conventions/api.md#errors-and-status-codes), [KI-05](../known-issues.md#ki-05) |
| No input validation; bad input becomes 0 or zero time | Agent assumes validation exists | [KI-06](../known-issues.md#ki-06) |
| Business rules exist only in the owner's head | Agent invents rules for stock and payment | [invariants/ticketing.md](../invariants/ticketing.md), [open-questions.md](../open-questions.md) |
| Hidden global state (`init()` singletons, `os.Getenv`, `LoadConfig` in `cmd` `init`) | Hard to unit-test; tests need a `.env` | [KI-10](../known-issues.md#ki-10) |
| **No tests at all** | No feedback loop; regressions go unnoticed | [KI-11](../known-issues.md#ki-11) |
| `go vet` fails, so `make lint` and `make test` fail | The agent's "verify" step is red from the start | [KI-08](../known-issues.md#ki-08) |

## Recommendations (in priority order)

1. Fix [KI-08](../known-issues.md#ki-08) so `make lint` and `make test` are green. An agent needs a green baseline.
2. Add a migration that fixes the schema drift ([KI-01](../known-issues.md#ki-01)–[KI-03](../known-issues.md#ki-03)) **before** building ticket categories & stock.
3. Add tests for the event service (mocked repository) and one integration test against Postgres. Stock work needs a concurrency test for [RULE-TKT-02](../invariants/ticketing.md).
4. Switch to sentinel not-found errors and `errors.Is` mapping ([KI-04](../known-issues.md#ki-04), [KI-05](../known-issues.md#ki-05)) before more domains copy the old pattern.
5. Answer the [open questions](../open-questions.md) that block ordering.
