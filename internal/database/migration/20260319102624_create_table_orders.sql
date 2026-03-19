-- +migrate Up
-- SQL untuk UP (Create/Alter table) ditulis di bawah baris ini
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT PRIMARY KEY,
    ticket_id BIGINT NOT NULL,
    bid_category_id BIGINT,
    order_id VARCHAR NOT NULL UNIQUE,
    email VARCHAR,
    phone VARCHAR,
    name VARCHAR,
    gender VARCHAR,
    bib_number VARCHAR,
    blood_type VARCHAR,
    emergency_call VARCHAR,
    payment_status VARCHAR NOT NULL DEFAULT 'pending',
    jersey_size VARCHAR,
    price_amount DECIMAL(12,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP,
    CONSTRAINT fk_orders_ticket_id FOREIGN KEY (ticket_id) REFERENCES tickets (id) ON DELETE RESTRICT
);

CREATE INDEX idx_orders_ticket_id ON orders(ticket_id);
CREATE INDEX idx_orders_order_id ON orders(order_id);
CREATE INDEX idx_orders_email ON orders(email);

-- +migrate Down
-- SQL untuk DOWN (Drop table/Rollback) ditulis di bawah baris ini
DROP INDEX IF EXISTS idx_orders_ticket_id;
DROP INDEX IF EXISTS idx_orders_order_id;
DROP INDEX IF EXISTS idx_orders_email;

DROP TABLE IF EXISTS orders;
