BEGIN;

ALTER TABLE history_items
    DROP COLUMN quantity_before,
    DROP COLUMN quantity_after,
    DROP COLUMN note;

END;