BEGIN;

DROP TABLE orders;

DROP TABLE payments;

DROP TABLE products_in_orders;

DROP TABLE statuses;

DROP EXTENSION IF EXISTS "uuid-ossp";

DROP SEQUENCE IF EXISTS seq_1;

COMMIT;
