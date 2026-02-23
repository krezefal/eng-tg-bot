-- =========================
-- UP migration
-- =========================
BEGIN;

ALTER TABLE user_dictionaries
    ADD COLUMN IF NOT EXISTS start_learning_at TIMESTAMPTZ NULL;

COMMIT;