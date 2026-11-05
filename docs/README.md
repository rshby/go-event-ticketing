# Docs — Single Source of Truth

This folder is the single source of truth (SSOT) for **what this service does and why**. Code tells you *how* something is built today. These docs record intent, business rules and decisions that cannot be recovered from the code.

When docs and code disagree:
- if the topic is a **business rule** (`invariants/`, `domain/`), the docs win and the code is a bug;
- if the topic is a **description of current behaviour** (`architecture/`), the code wins and the doc must be fixed in the same change.

> Citations use `path:line` plus a symbol name. Line numbers drift, so treat the **symbol** as authoritative.

## Folders

| Path | What it holds | Read it before… |
|------|---------------|-----------------|
| [domain/glossary.md](domain/glossary.md) | Business vocabulary (event, ticket tier, category, order, bib…) | anything — read this first |
| [architecture/](architecture/README.md) | How the service is built: layers, infra, per-domain design | adding or changing any feature |
| [invariants/](invariants/README.md) | Rules that must never be broken (`RULE-XXX-NN`) | touching tickets, orders, IDs, cache, migrations |
| [conventions/](conventions/) | How **new** code must be written (Go, API, DB, observability) | writing code |
| [reqspecs/](reqspecs/README.md) | Requirement specs for upcoming features | starting a new feature |
| [implspecs/](implspecs/README.md) | Dated implementation plans | implementing an approved feature |
| [known-issues.md](known-issues.md) | Defects and drift found in the current code | copying an existing pattern |
| [open-questions.md](open-questions.md) | Business decisions not yet made — **do not invent answers** | designing tickets / orders / payment |
| [superpowers/specs/](superpowers/specs/) | Design documents produced by brainstorming sessions | understanding why a doc set exists |

## Conventions

- [conventions/golang.md](conventions/golang.md) — layering, constructors, error flow, tracing/logging boilerplate
- [conventions/api.md](conventions/api.md) — routes, response envelope, status codes, JSON casing, validation
- [conventions/database.md](conventions/database.md) — migrations, GORM, IDs, money
- [conventions/observability.md](conventions/observability.md) — spans, middleware, log fields

## Rules for keeping docs true

1. A code change that alters behaviour described here updates the doc **in the same change**.
2. Anything not built yet is labelled **Planned**. Anything undecided goes in [open-questions.md](open-questions.md) and is labelled **Open question** where it is referenced.
3. When an open question is answered, move the answer into the relevant doc (glossary, architecture or invariants) and delete it from `open-questions.md`.
4. When a known issue is fixed, delete its entry from [known-issues.md](known-issues.md).
