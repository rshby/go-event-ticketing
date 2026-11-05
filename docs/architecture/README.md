# Architecture

Technical design for each part of the service. Each domain file covers its data model, endpoints, infrastructure (tables, Redis keys) and behaviour, and marks clearly what is **Built** and what is **Planned**.

| File | Covers |
|------|--------|
| [overview.md](overview.md) | Whole service: layers, request flow, composition root, infrastructure, config, IDs, time |
| [event.md](event.md) | Event domain — **Built** |
| [ticketing.md](ticketing.md) | Ticket tiers, categories, orders and stock — schema exists, logic **Planned** |
| [review.md](review.md) | Assessment of the architecture for AI-assisted development |

## Adding a new domain

1. Write a reqspec in [../reqspecs/](../reqspecs/README.md) and get it approved.
2. Add or extend `docs/architecture/<domain>.md` following the shape of [event.md](event.md).
3. If the domain has rules that must never break, add `docs/invariants/<domain>.md`.
4. Write an implspec in [../implspecs/](../implspecs/README.md), then implement.
5. After merging, update the architecture doc so it describes what was built.
