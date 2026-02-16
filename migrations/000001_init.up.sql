-- =========================
-- UP migration
-- =========================
BEGIN;

-- Extensions
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Enums
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'dictionary_mode') THEN
CREATE TYPE dictionary_mode AS ENUM ('random_pool', 'on_schedule');
END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_word_status') THEN
CREATE TYPE user_word_status AS ENUM ('learning', 'blocked');
END IF;
END$$;

-- users
CREATE TABLE IF NOT EXISTS users (
    tg_id                BIGINT PRIMARY KEY,
    active_dictionary_id UUID NULL      -- активный словарь пользака
);

-- dictionaries
CREATE TABLE IF NOT EXISTS dictionaries (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title       VARCHAR(25) NOT NULL,
    description VARCHAR(50) NOT NULL DEFAULT '',
    mode        dictionary_mode NOT NULL,
    author      VARCHAR(50) NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- FK на активный словарь - если словарь удалили, сбросить активный
ALTER TABLE users
    ADD CONSTRAINT fk_users_active_dictionary
        FOREIGN KEY (active_dictionary_id)
            REFERENCES dictionaries(id)
            ON DELETE SET NULL;

-- dictionary_schedule_batch
CREATE TABLE IF NOT EXISTS dictionary_schedule_batch (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dictionary_id UUID NOT NULL REFERENCES dictionaries(id) ON DELETE CASCADE,
    delay_days    INT NOT NULL CHECK (delay_days >= 0)
);

-- composite FK из dictionary_words(dictionary_id, batch_id)
ALTER TABLE dictionary_schedule_batch
    ADD CONSTRAINT uq_dictionary_schedule_batch_id_dictionary_id
        UNIQUE (dictionary_id, id);

CREATE INDEX IF NOT EXISTS idx_dictionary_schedule_batch_dictionary_id
    ON dictionary_schedule_batch(dictionary_id);

-- dictionary_words
CREATE TABLE IF NOT EXISTS dictionary_words (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dictionary_id  UUID NOT NULL REFERENCES dictionaries(id) ON DELETE CASCADE,
    batch_id       UUID REFERENCES dictionary_schedule_batch(id) ON DELETE SET NULL,
    spelling       VARCHAR(25) NOT NULL,
    transcription  VARCHAR(25) NOT NULL DEFAULT '',
    audio          TEXT NOT NULL DEFAULT '',   -- mp3 ref: url/path/key;
    ru_translation VARCHAR(25) NOT NULL DEFAULT ''
);

-- composite FK: batch должен принадлежать тому же dictionary_id
ALTER TABLE dictionary_words
    ADD CONSTRAINT fk_dictionary_words_dictionary_batch_pair
        FOREIGN KEY (dictionary_id, batch_id)
            REFERENCES dictionary_schedule_batch(dictionary_id, id);

-- уникальность слова в пределах словаря
ALTER TABLE dictionary_words
    ADD CONSTRAINT uq_dictionary_words_dictionary_spelling
        UNIQUE (dictionary_id, spelling);

CREATE INDEX IF NOT EXISTS idx_dictionary_words_dictionary_id
    ON dictionary_words(dictionary_id);

CREATE INDEX IF NOT EXISTS idx_dictionary_words_dictionary_id_batch_id
    ON dictionary_words(dictionary_id, batch_id);

-- user_dictionaries
CREATE TABLE IF NOT EXISTS user_dictionaries (
    user_id       BIGINT NOT NULL REFERENCES users(tg_id) ON DELETE CASCADE,
    dictionary_id UUID   NOT NULL REFERENCES dictionaries(id) ON DELETE CASCADE,
    subscribed_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, dictionary_id)
);

CREATE INDEX IF NOT EXISTS idx_user_dictionaries_dictionary_id
    ON user_dictionaries(dictionary_id);

-- user_words_state
CREATE TABLE IF NOT EXISTS user_words_state (
    user_id        BIGINT NOT NULL REFERENCES users(tg_id) ON DELETE CASCADE,
    dict_word_id   UUID   NOT NULL REFERENCES dictionary_words(id) ON DELETE CASCADE,
    status         user_word_status NOT NULL DEFAULT 'learning',

    ef             REAL NOT NULL DEFAULT 2.5,
    interval_days  INT  NOT NULL DEFAULT 0 CHECK (interval_days >= 0),
    repetition     INT  NOT NULL DEFAULT 0 CHECK (repetition >= 0),

    last_result    INT NULL CHECK (last_result BETWEEN 0 AND 5),
    last_review_at TIMESTAMPTZ NULL,
    next_review_at TIMESTAMPTZ NULL,

    PRIMARY KEY (user_id, dict_word_id)
);

CREATE INDEX IF NOT EXISTS idx_user_words_state_user_status
    ON user_words_state(user_id, status);

CREATE INDEX IF NOT EXISTS idx_user_words_state_user_next_review
    ON user_words_state(user_id, next_review_at);

COMMIT;
