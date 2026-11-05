# Glossary

The service sells **registrations for running races**. All terms below are business terms. The code name follows in brackets.

| Term | Code | Meaning |
|------|------|---------|
| **Event** | `entity.Event`, table `events` | One race (e.g. "Jakarta Marathon 2026"). Has a name, banner, location and a start/end time. |
| **Ticket** / **ticket tier** | `entity.Ticket`, table `tickets` | A **sales tier** of an event, e.g. *Early Bird* or *Normal*. Each tier has its own price, sale window (`start_time`–`end_time`) and quota. One event has many tiers. A ticket is **not** a single seat or a single runner. |
| **Category** | `entity.TicketCategory`, table `ticket_categories` | The **race distance** the runner chooses, e.g. 5K, 10K, Half Marathon. A label on the order. **Has no stock of its own.** |
| **Quota** | `tickets.max_quota`, `tickets.available_quota` | Stock is counted **per ticket tier only**. `max_quota` = total slots in the tier; `available_quota` = slots not yet reserved or sold. |
| **Order** | `entity.Order`, table `orders` | One runner's registration for one ticket tier and one category. Holds runner data (name, email, phone, gender, blood type, emergency contact, jersey size, bib number) and the price paid. |
| **Order ID** | `orders.order_id` | Public, human-facing order reference (unique string). Different from the internal numeric `orders.id`. Format: **Open question**. |
| **Payment status** | `orders.payment_status` | State of the order's payment. Starts as `pending`. Other states: **Open question**. |
| **Reservation** | — | Holding one slot of a tier's quota for a `pending` order until it is paid or expires. |
| **Expiry** | `orders.expired_at` | The time a `pending` order stops holding its reservation. After this the slot goes back to `available_quota`. |
| **Bib number** | `orders.bib_number` | The race number printed on the runner's bib. How it is assigned: **Open question**. |
| **Jersey size** | `orders.jersey_size` | Race shirt size chosen by the runner. |

## Relationships

```
Event 1 ──< Ticket (tier)  1 ──< Order >── 1 Category (distance)
             quota lives here        runner data lives here
```

- An event has many ticket tiers (`tickets.event_id` → `events.id`, `ON DELETE CASCADE`).
- An order belongs to exactly one ticket tier (`orders.ticket_id` → `tickets.id`, `ON DELETE RESTRICT`) and one category.
- Whether categories are a global list or belong to an event: **Open question** (see [../open-questions.md](../open-questions.md)).

See [../architecture/ticketing.md](../architecture/ticketing.md) for the full data model and the stock lifecycle.
