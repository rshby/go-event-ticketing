# Database Conventions

## Migrations

1. `go run . create-migration-file <snake_case_name>` creates `internal/database/migration/<YYYYMMDDhhmmss>_<name>.sql` containing `-- +migrate Up` / `-- +migrate Down` sections.
2. Write the `Up` SQL **and** a working `Down`.
3. `go run . db-migrate` applies pending `Up` migrations (table `migrations` tracks what has run). There is no CLI command for `Down` yet.
4. Never edit a migration that may have been applied. Fix forward with a new file ([RULE-PLT-02](../invariants/platform.md)).
5. No `AutoMigrate`.

SQL style for new files: uppercase keywords, `IF NOT EXISTS` / `IF EXISTS`, explicit constraint names (`fk_<table>_<column>`, `idx_<table>_<column>`, `chk_<table>_<rule>`).

## Tables

- Table names: plural snake_case (`tickets`, `orders`, `ticket_categories`).
- Every table has `id BIGINT PRIMARY KEY`, `created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`, `updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP`.
- Foreign keys are declared, and every FK column is indexed.
- Use DB constraints for invariants that the DB can enforce (`NOT NULL`, `UNIQUE`, `CHECK`).

## GORM entities

- One struct per table in `internal/entity/<domain>.go`, with `TableName()` returning the exact SQL name.
- Every field has `gorm:"column:<name>"`, and the tags match the migration ([RULE-PLT-02](../invariants/platform.md)). After writing a migration, compare it with the struct.
- Queries: `r.db.WithContext(ctx).Model(&entity.X{})…` — always `WithContext(ctx)`, so spans and cancellation work.
- Conditional updates use `Where(...).UpdateColumn(...)` / `Updates(...)` and check `RowsAffected`.
- `SkipDefaultTransaction: true` is set, so a multi-statement write needs an explicit `Transaction`.

## IDs

The repository assigns `ID` with `helper.GenerateSonyFlakeID()` before insert, and returns the error if generation fails. DTO mappers don't generate IDs ([RULE-PLT-01](../invariants/platform.md), [KI-14](../known-issues.md#ki-14)).

## Money

SQL: `DECIMAL(12,2)`. Go: `decimal.Decimal` from `github.com/shopspring/decimal`, which supports GORM scanning. Never `float64` ([KI-07](../known-issues.md#ki-07)).

## Time

Columns are `TIMESTAMP` without time zone, holding Asia/Jakarta wall-clock time. Use `helper.GetNowTime()` in Go. Within one statement, comparing with `now()` in SQL is fine, as long as `DB_TIMEZONE` is `Asia/Jakarta`.
