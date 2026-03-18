-- +migrate Up
-- SQL untuk UP (Create/Alter table) ditulis di bawah baris ini
CREATE TABLE IF Not Exists events (
    id BIGINT PRIMARY KEY,
    name VARCHAR NOT NULL,
    banner_image VARCHAR,
    location VARCHAR,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR
);

-- +migrate Down
-- SQL untuk DOWN (Drop table/Rollback) ditulis di bawah baris ini
DROP TABLE IF EXISTS events;
