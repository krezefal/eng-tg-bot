-- =========================
-- DOWN migration
-- =========================
BEGIN;

DROP TRIGGER IF EXISTS trg_cleanup_user_progress_on_unsubscribe ON user_dictionaries;
DROP FUNCTION IF EXISTS cleanup_user_progress_on_unsubscribe();

COMMIT;
