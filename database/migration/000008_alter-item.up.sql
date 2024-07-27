BEGIN;

ALTER TABLE items
    ADD COLUMN keywords TEXT;

CREATE INDEX idx_items_keywords ON items (keywords);

COMMIT;
