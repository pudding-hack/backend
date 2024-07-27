BEGIN;

ALTER TABLE history_items
    ADD COLUMN quantity_before bigint,
    ADD COLUMN quantity_after bigint,
    ADD COLUMN note TEXT;

COMMIT;
