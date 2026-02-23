-- =========================
-- DOWN migration
-- =========================
BEGIN;

ALTER TABLE user_dictionaries
DROP COLUMN IF EXISTS start_learning_at;

COMMIT;