# Implementation Specs

An implspec is the step-by-step build plan for an approved [reqspec](../reqspecs/README.md). It must be precise enough that an agent can execute it without asking questions.

- File name: `docs/implspecs/YYYY-MM-DD-<domain>-<feature>.md`.
- Every step names the files it touches and how it is verified.
- It cites the invariants it must satisfy and the known issues it fixes or works around.
- When the work is merged, update [../architecture/](../architecture/README.md) so it describes what was actually built. The implspec stays as a historical record.

## Template

```markdown
# <Domain> — <Feature> (implementation)

**Reqspec:** ../reqspecs/<file>.md
**Date:** YYYY-MM-DD
**Status:** Draft | Approved | Done

## Summary
<one paragraph: what changes and the approach>

## Invariants
- RULE-TKT-02 → enforced by <mechanism>, proven by <test>

## Steps
### 1. Migration <name>
Files: internal/database/migration/<ts>_<name>.sql
- …
Verify: `go run . db-migrate` then `\d <table>`

### 2. Entity + interfaces
Files: internal/entity/<domain>.go
…

### N. Docs
Files: docs/architecture/<domain>.md, docs/known-issues.md (remove fixed KIs)

## Test plan
- Unit: …
- Integration / concurrency: …

## Verification
`make lint && make test && make build`, plus manual curl calls: …
```
