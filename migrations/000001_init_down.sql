-- =========================
-- DOWN migration
-- =========================
BEGIN;

DROP TABLE IF EXISTS user_words_state;
DROP TABLE IF EXISTS user_dictionaries;
DROP TABLE IF EXISTS dictionary_words;
DROP TABLE IF EXISTS dictionary_schedule_batch;
DROP TABLE IF EXISTS dictionaries;
DROP TABLE IF EXISTS users;

DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_word_status') THEN
DROP TYPE user_word_status;
END IF;

  IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'dictionary_mode') THEN
DROP TYPE dictionary_mode;
END IF;
END$$;

COMMIT;
