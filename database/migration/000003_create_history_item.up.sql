BEGIN;

CREATE TABLE history_items (
    id serial primary key,
    item_id bigint,
    qty bigint default 0,
    end_qty bigint default 0,
    start_qty bigint default 0,
    type_id int default 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    created_by VARCHAR(255) DEFAULT 'system',
    updated_by VARCHAR(255) DEFAULT 'system',
    deleted_by VARCHAR(255)
);

COMMIT;