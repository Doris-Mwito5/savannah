-- +goose Up
CREATE TYPE ORDER_STATUS AS ENUM ('pending', 'paid', 'cancelled', 'returned');
CREATE TYPE ORDER_SOURCE AS ENUM('offline', 'online');
CREATE TYPE PAYMENT_METHOD AS ENUM('cash', 'mpesa', 'card');

CREATE TABLE orders (
    id                  BIGSERIAL           PRIMARY KEY,
    reference_number    VARCHAR(40)         NOT NULL,
    phone_number        VARCHAR(20)         NOT NULL,
    order_status        ORDER_STATUS        NOT NULL DEFAULT 'pending',
    order_source        ORDER_SOURCE        NOT NULL DEFAULT 'offline',
    payment_method      PAYMENT_METHOD      NOT NULL DEFAULT 'cash',
    customer_id         BIGINT              NOT NULL REFERENCES customers(id),
    shop_id             VARCHAR(40),
    total_items         INTEGER             NOT NULL DEFAULT 0,
    total_amount        DECIMAL(10,2)       NOT NULL DEFAULT 0.00,
    discount            DECIMAL(10, 2)      DEFAULT 0.00,
    created_at          TIMESTAMPTZ         NOT NULL DEFAULT clock_timestamp(),
    updated_at          TIMESTAMPTZ         NOT NULL DEFAULT clock_timestamp()
);

CREATE UNIQUE INDEX orders_reference_number_uniq_idx ON orders(reference_number);
CREATE INDEX orders_phone_number_idx ON orders(phone_number);
CREATE INDEX orders_payment_method_idx ON orders(payment_method);
CREATE INDEX orders_order_status_idx ON orders(order_status);
CREATE INDEX orders_order_source_idx ON orders(order_source);
CREATE INDEX orders_customer_id_idx ON orders(customer_id);

CREATE TABLE order_items(
    id              BIGSERIAL       PRIMARY KEY,
    order_id        BIGINT          NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id      BIGINT          NOT NULL REFERENCES products(id),
    unit_price      DECIMAL(10, 2)  NOT NULL,
    quantity        INTEGER         NOT NULL CHECK (quantity > 0),
    total_amount    DECIMAL(10, 2)  NOT NULL, -- quantity * unit_price
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT clock_timestamp(),
    updated_at          TIMESTAMPTZ         NOT NULL DEFAULT clock_timestamp()
);

CREATE INDEX order_items_order_id_idx ON order_items(order_id);
CREATE INDEX order_items_product_id_idx ON order_items(product_id);

-- +goose Down
DROP INDEX IF EXISTS order_items_product_id_idx;
DROP INDEX IF EXISTS order_items_order_id_idx;
DROP TABLE IF EXISTS order_items;

DROP INDEX IF EXISTS orders_customer_id_idx;
DROP INDEX IF EXISTS orders_order_source_idx;
DROP INDEX IF EXISTS orders_order_status_idx;
DROP INDEX IF EXISTS orders_payment_method_idx;
DROP INDEX IF EXISTS orders_phone_number_idx;
DROP INDEX IF EXISTS orders_reference_number_uniq_idx;
DROP TABLE IF EXISTS orders;

DROP TYPE IF EXISTS PAYMENT_METHOD;
DROP TYPE IF EXISTS ORDER_SOURCE;
DROP TYPE IF EXISTS ORDER_STATUS;