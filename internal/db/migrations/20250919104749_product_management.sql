-- +goose Up
CREATE TABLE categories (
    id                    BIGSERIAL               PRIMARY KEY,
    name                  VARCHAR(255)         NOT NULL,
    parent_id             BIGINT               REFERENCES categories(id),
    shop_id               VARCHAR(180),
    created_at            TIMESTAMPTZ          NOT NULL DEFAULT clock_timestamp(),
    updated_at            TIMESTAMPTZ          NOT NULL DEFAULT clock_timestamp()
);

CREATE INDEX categories_name_idx ON categories(name);

CREATE TYPE PRODUCT_TYPE AS ENUM('goods', 'services');

CREATE TABLE products(
    id                   BIGSERIAL                PRIMARY KEY,
    name                 VARCHAR(255)          NOT NULL,
    description          TEXT,
    wholesale_price      DECIMAL(10, 4)        NOT NULL,
    retail_price         DECIMAL(10, 4)        NOT NULL,
    category_id          BIGINT                NOT NULL REFERENCES categories(id),
    product_image        TEXT,
    stock                INTEGER,
    product_type         PRODUCT_TYPE          NOT NULL DEFAULT 'goods',
    created_at           TIMESTAMPTZ           NOT NULL DEFAULT clock_timestamp(),
    updated_at           TIMESTAMPTZ           NOT NULL DEFAULT clock_timestamp()
);

-- Regular index, NOT unique - multiple products can belong to same category
CREATE INDEX products_category_id_idx ON products(category_id);
CREATE INDEX products_name_idx ON products(name);

-- +goose Down
DROP INDEX IF EXISTS products_name_idx;
DROP INDEX IF EXISTS products_category_id_idx;
DROP TABLE IF EXISTS products;
DROP TYPE IF EXISTS PRODUCT_TYPE;

DROP INDEX IF EXISTS categories_name_idx;
DROP TABLE IF EXISTS categories;