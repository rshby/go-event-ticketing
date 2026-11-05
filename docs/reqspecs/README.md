# Requirement Specs

A reqspec records **what** a feature must do and **why**, before any design or code. It is where human decisions are written down, so an agent doesn't have to guess.

- File name: `docs/reqspecs/<domain>-<feature>.md` (e.g. `ticketing-categories-and-stock.md`).
- Write one before any new feature. Bug fixes and refactors don't need one.
- Link the [open questions](../open-questions.md) the feature depends on. A reqspec with unresolved blocking questions isn't ready.
- After approval, write an implspec ([../implspecs/](../implspecs/README.md)).

## Template

```markdown
# <Domain> — <Feature>

**Status:** Draft | Approved
**Owner:** <name>
**Date:** YYYY-MM-DD

## Why
<problem, who has it, what happens today>

## Scope
In: …
Out (non-goals): …

## Functional requirements
- **FR-01:** <the system must …>
- **FR-02:** …

## Acceptance criteria
- **AC-01 (FR-01):** Given <state>, when <action>, then <result>.
- **AC-02 (FR-01):** Given … (include sold-out, concurrency and error cases)

## Invariants touched
- RULE-TKT-02 — <how this feature keeps it>

## Open questions
- OQ-xx — <blocking? yes/no>
```
