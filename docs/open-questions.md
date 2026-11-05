# Open Questions

Business decisions that **only the product owner can make**. They are not answered anywhere in code or docs.

**For agents:** do not invent answers. If a task depends on one of these, stop and ask, or implement only the part that doesn't depend on it, and say so. When an answer is given, write it into the relevant doc and delete the entry here.

| ID | Blocks |
|----|--------|
| [OQ-01](#oq-01) | order states, quota release paths |
| [OQ-02](#oq-02) | expiry job |
| [OQ-03](#oq-03) | payment integration |
| [OQ-04](#oq-04) | bib assignment |
| [OQ-05](#oq-05) | category model |
| [OQ-06](#oq-06) | order limits |
| [OQ-07](#oq-07) | order reference |
| [OQ-08](#oq-08) | admin / auth |
| [OQ-09](#oq-09) | editing tier quota |
| [OQ-10](#oq-10) | deleting tiers / categories |
| [OQ-11](#oq-11) | scope of "categories & stock" |
| [OQ-12](#oq-12) | tiers without a sale window |

<a id="oq-01"></a>

#### OQ-01: Full set of order states
Known: `pending` (default), and working names `paid`, `expired`. Unknown: cancelled by user? payment failed? refunded? Does each one return the slot? ([RULE-TKT-04](invariants/ticketing.md), [RULE-TKT-05](invariants/ticketing.md))

<a id="oq-02"></a>

#### OQ-02: Reservation length and expiry mechanism
How long after creation is `expired_at` (15 min? 1 h? 24 h)? Is it the same for every tier? What runs expiry: a cron/ticker in the service, a separate command, Redis key expiry events, or a lazy check on read?

<a id="oq-03"></a>

#### OQ-03: Payment
Which payment provider? How is payment confirmed (webhook, polling, manual by admin)? What happens when a payment arrives after the order expired: refund, or re-reserve if stock exists?

<a id="oq-04"></a>

#### OQ-04: Bib numbers
Who assigns bib numbers, and when (on payment, at a cut-off date, manually)? Is the format per category (e.g. 10K = 1xxxx)? Must they be unique per event?

<a id="oq-05"></a>

#### OQ-05: Category scope
Are categories a global list shared by all events, or does each event define its own (`ticket_categories.event_id`)? Can a ticket tier be limited to certain categories (e.g. Early Bird for 10K only)?

<a id="oq-06"></a>

#### OQ-06: Order limits
One order = one runner? Max orders per email/phone per event? Can the same email register for two categories?

<a id="oq-07"></a>

#### OQ-07: `order_id` format
What does the public order reference look like (e.g. `ORD-20261201-XXXX`)? Who generates it?

<a id="oq-08"></a>

#### OQ-08: Who can create events and tiers
There is no authentication; `created_by` is hard-coded to `"admin"`. Are admin endpoints protected? Is there a runner account, or are orders anonymous by email?

<a id="oq-09"></a>

#### OQ-09: Editing a tier's `max_quota`
When an admin changes `max_quota`, how does `available_quota` change (shift by the same difference)? Can `max_quota` drop below the number of slots already reserved or sold, and if so, is it rejected?

<a id="oq-10"></a>

#### OQ-10: Deleting tiers and categories that have orders
Should this be blocked (409), soft-deleted, or allowed? Today an FK `RESTRICT` makes it a 500. What `ON DELETE` action should the new `orders.ticket_category_id` FK use?

<a id="oq-11"></a>

#### OQ-11: Scope of the "ticket categories & stock" feature
Is it tier + category CRUD only (stock fields and the atomic decrement ready, but no orders), or does it include create-order, pay and expiry? Ordering is also blocked by OQ-01, OQ-02, OQ-03, OQ-06 and OQ-07.

<a id="oq-12"></a>

#### OQ-12: Tiers without a sale window
`tickets.start_time` and `end_time` are nullable. Must every tier have both (validation + `NOT NULL` migration), or does NULL mean "always on sale"? With `BETWEEN`, NULL today means "never on sale".
