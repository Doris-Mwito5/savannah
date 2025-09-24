-- +goose Up
CREATE TYPE CUSTOMER_TYPE AS ENUM ('individual', 'business');

CREATE TABLE customers (
    id                  BIGSERIAL              PRIMARY KEY,
    name                VARCHAR(50)         NOT NULL,
    email               VARCHAR(40)         NOT NULL,
    phone_number        VARCHAR(20)         UNIQUE NOT NULL,
    customer_type       CUSTOMER_TYPE       NOT NULL DEFAULT 'individual',
    shop_id             VARCHAR(155),
    created_at          TIMESTAMPTZ         NOT NULL DEFAULT clock_timestamp(),
    updated_at          TIMESTAMPTZ         NOT NULL DEFAULT clock_timestamp()
);

CREATE UNIQUE INDEX customers_name_idx ON customers(name);
CREATE UNIQUE INDEX customers_email_idx ON customers(email);

-- +goose Down
DROP INDEX IF EXISTS customers_email_idx;
DROP INDEX IF EXISTS customers_name_idx;

DROP TABLE IF EXISTS customers;
DROP TYPE IF EXISTS CUSTOMER_TYPE;
