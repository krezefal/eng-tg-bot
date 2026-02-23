-- =========================
-- UP migration
-- =========================
BEGIN;

CREATE OR REPLACE FUNCTION cleanup_user_progress_on_unsubscribe()
RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM user_words_state uws
    USING dictionary_words dw
    WHERE uws.user_id = OLD.user_id
      AND uws.dict_word_id = dw.id
      AND dw.dictionary_id = OLD.dictionary_id;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_cleanup_user_progress_on_unsubscribe ON user_dictionaries;

CREATE TRIGGER trg_cleanup_user_progress_on_unsubscribe
AFTER DELETE ON user_dictionaries
FOR EACH ROW
EXECUTE FUNCTION cleanup_user_progress_on_unsubscribe();

COMMIT;
