BEGIN;

ALTER TABLE items
    DROP COLUMN keywords;

DROP INDEX idx_items_keywords;

END;