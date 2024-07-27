BEGIN;

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    item_code varchar(15),
    item_name varchar(255),
    qty bigint default 0,
    unit bigint,
    price DOUBLE PRECISION DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by VARCHAR(255) DEFAULT 'system',
    updated_by VARCHAR(255) DEFAULT 'system',
    deleted_by VARCHAR(255)
);

COMMIT;