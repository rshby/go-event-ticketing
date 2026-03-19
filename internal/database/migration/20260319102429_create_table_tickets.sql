-- +migrate Up
-- SQL untuk UP (Create/Alter table) ditulis di bawah baris ini
create table if not exists tickets (
    id BIGINT PRIMARY KEY,
    event_id BIGINT,
    name VARCHAR NOT NULL,
    banner_image VARCHAR,
    max_quota INTEGER NOT NULL,
    available_quota INTEGER NOT NULL,
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    price DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_tickets_event_id FOREIGN KEY (event_id) REFERENCES events (id) ON DELETE CASCADE
);

CREATE INDEX idx_tickets_event_id ON tickets(event_id);

-- +migrate Down
-- SQL untuk DOWN (Drop table/Rollback) ditulis di bawah baris ini
drop index if exists idx_tickets_event_id;
drop table if exists tickets;
