-- +migrate Up
-- SQL untuk UP (Create/Alter table) ditulis di bawah baris ini
CREATE TABLE IF NOT EXISTS event_category (
    id BIGINT PRIMARY KEY,
    name VARCHAR NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- +migrate Down
-- SQL untuk DOWN (Drop table/Rollback) ditulis di bawah baris ini
DROP TABLE IF EXISTS event_category;
