BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SEQUENCE IF NOT EXISTS seq_1
INCREMENT 1
START 100000
MINVALUE  100000
MAXVALUE 999999
CYCLE;

CREATE TABLE orders (
    id VARCHAR(36) PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(), 
    short INTEGER DEFAULT nextval('seq_1'), 
    date TIMESTAMP, 
    status_id SMALLINT
);

CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY NOT NULL, 
    total REAL, 
    change REAL, 
    order_id VARCHAR(36)
);

CREATE TABLE products_in_orders (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    uuid VARCHAR(36),
    num INTEGER,
    price_per_one REAL, 
    order_id VARCHAR(36)
);

CREATE TABLE statuses (
    id SMALLINT PRIMARY KEY NOT NULL,
    name VARCHAR(40)
);

COMMIT;