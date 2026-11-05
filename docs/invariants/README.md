# Invariants

Rules that must **never** be broken, whatever the feature. Each rule is something a reviewer (human or AI) can check in a diff.

| File | Prefix | Scope |
|------|--------|-------|
| [platform.md](platform.md) | `RULE-PLT` | IDs, schema, cache, errors — applies to every domain |
| [ticketing.md](ticketing.md) | `RULE-TKT` | Ticket quota, orders, reservation and expiry |

## Format

```
#### RULE-<DOMAIN>-<NN>: <short title> `[SEVERITY]` `[Planned]`?
**Invariant:** the rule, stated so it can be checked.
**Why:** the business or technical reason.
**Check for:** what to look for in code or a diff.
**Violation looks like:** a concrete wrong example.
```

| Severity | Meaning |
|----------|---------|
| `CRITICAL` | Breaking it loses money, oversells, corrupts data or leaks personal data |
| `HIGH` | Breaking it causes wrong behaviour that users notice |
| `MEDIUM` | Breaking it causes inconsistency or makes future bugs likely |

`[Planned]` means no code enforces the rule yet, because the feature isn't built. It is still binding: the first implementation must satisfy it.

## Changing a rule

Rules record decisions made by the product owner. An agent must not change or relax a rule on its own. Propose the change, get approval, then edit the rule and the code together.
