# Ticketing — Invariants

**Responsibility:** ticket tiers (sales tiers with quota), race categories, and orders (runner registrations).
**Key paths:** `internal/entity/ticket.go`, `internal/entity/ticket_category.go`, `internal/entity/order.go`, and future `internal/{repository,service,controller}/{ticket,order}*.go`
**Design:** [../architecture/ticketing.md](../architecture/ticketing.md) · **Vocabulary:** [../domain/glossary.md](../domain/glossary.md)

All rules here are `[Planned]`: the ordering code does not exist yet. They are decisions, not suggestions.

---

#### RULE-TKT-01: Quota stays within bounds `[CRITICAL]` `[Planned]`
**Invariant:** for every ticket tier, `0 ≤ available_quota ≤ max_quota`, always.
**Why:** a negative value means oversold race slots. A value above max means slots were released twice.
**Check for:** any write to `available_quota`. Prefer adding a DB `CHECK (available_quota BETWEEN 0 AND max_quota)` in the migration that builds this feature.
**Violation looks like:** `UPDATE tickets SET available_quota = available_quota - 1 WHERE id = ?` with no `available_quota > 0` guard.

#### RULE-TKT-02: Reservation is one atomic, conditional decrement `[CRITICAL]` `[Planned]`
**Invariant:** taking a slot is a single SQL `UPDATE … SET available_quota = available_quota - 1 WHERE id = ? AND available_quota > 0 AND now() BETWEEN start_time AND end_time`, and `RowsAffected == 1` decides success. It runs in the **same transaction** as the order insert.
**Why:** two runners buying the last slot at the same moment must not both succeed.
**Check for:** a read of `available_quota` into Go followed by a write; `GORM Save` on a `Ticket` loaded earlier; the decrement and the insert in separate transactions; stock decisions made from Redis.
**Violation looks like:** `if ticket.AvailableQuota > 0 { ticket.AvailableQuota--; db.Save(&ticket) }`.
**Reporting the reason:** when `RowsAffected == 0`, the transaction is rolled back. A plain read of the tier **after** the failure picks the error: missing → `ErrTicketNotFound` (404), outside window → `ErrTicketNotOnSale` (409), otherwise `ErrTicketSoldOut` (409). That read only chooses the message; it never retries or decides stock.

#### RULE-TKT-03: A pending order holds exactly one slot `[CRITICAL]` `[Planned]`
**Invariant:** creating an order in `pending` state reserves one slot from its `ticket_id` tier at creation time. No slot → no order (reject; the order row is not inserted).
**Why:** decided by the product owner: reserve on order, not on payment, so a runner who starts paying is not beaten by someone else.
**Check for:** order creation paths that insert the order without decrementing, or decrement only on payment.
**Violation looks like:** quota decremented in the payment-confirmation handler.

#### RULE-TKT-04: A slot is returned exactly once when an order leaves `pending` unpaid `[CRITICAL]` `[Planned]`
**Invariant:** when a `pending` order ends without payment, its slot is added back to `available_quota` **exactly once**. Today the only such path is **expiry** (`now ≥ expired_at`). Cancellation and failed payment, if added ([OQ-01](../open-questions.md#oq-01)), follow the same rule.
**Why:** a lost release locks the slot forever; a double release oversells.
**Check for:** the state change and the increment must be one transaction, guarded by the current state: `UPDATE orders SET payment_status='expired' WHERE id=? AND payment_status='pending'`, and increment only if that returned `RowsAffected == 1`. Expiry jobs must be safe to run twice or concurrently.
**Violation looks like:** an expiry job that increments quota for every order with `expired_at < now()` without first moving the order out of `pending`.

#### RULE-TKT-05: Paid orders never release their slot automatically `[HIGH]` `[Planned]`
**Invariant:** an order that has been paid is never expired, and its slot is never returned by the expiry path.
**Why:** a race between payment confirmation and expiry must not give away a paid runner's slot.
**Check for:** payment confirmation must also be guarded (`WHERE payment_status='pending'`). If expiry won the race, payment confirmation must fail visibly, not silently reopen the order. How to handle a payment that arrives after expiry: [OQ-03](../open-questions.md#oq-03).
**Violation looks like:** `UPDATE orders SET payment_status='paid' WHERE id=?` with no state guard.

#### RULE-TKT-06: Orders only during the tier's sale window `[HIGH]` `[Planned]`
**Invariant:** an order for tier T can be created only when `T.start_time ≤ now ≤ T.end_time`, measured in the business time zone (Asia/Jakarta).
**Check for:** order creation that checks only quota.
**Violation looks like:** an Early Bird order accepted after the Early Bird window closed.

#### RULE-TKT-07: Order price is a snapshot `[HIGH]` `[Planned]`
**Invariant:** `orders.price_amount` is copied from `tickets.price` when the order is created, and never recalculated. Later price changes to the tier do not affect existing orders. The client never sends the price.
**Check for:** a price taken from the request body; joins to `tickets.price` when showing an order's price.

#### RULE-TKT-08: Category is a label, not stock `[MEDIUM]` `[Planned]`
**Invariant:** choosing a category never changes any quota. Categories have no quota columns.
**Why:** decided: quota is per ticket tier only.
**Check for:** new `quota` columns on `ticket_categories`, or decrements keyed by category.

#### RULE-TKT-09: Runner personal data is not logged or traced `[CRITICAL]` `[Planned]`
**Invariant:** phone, email, blood type, emergency contact and name from orders do not appear in logs or span attributes in full.
**Why:** this is personal data, including health data.
**Check for:** `helper.Dump(&request)` of order DTOs in log fields; `TraceMiddleware` body capture on order routes ([KI-22](../known-issues.md#ki-22)).

#### RULE-TKT-10: Quota changes invalidate cached ticket data `[MEDIUM]` `[Planned]`
**Invariant:** if ticket tiers are cached, every reservation, release and tier edit deletes the affected cache keys. Remaining quota shown to users may come from cache. **Reservation decisions never do** ([RULE-PLT-03](platform.md)).
