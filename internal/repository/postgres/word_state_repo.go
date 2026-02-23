package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

type WordsStateRepo struct {
	db     *sql.DB
	logger *zerolog.Logger
}

func NewWordsStateRepo(db *sql.DB, parentLogger *zerolog.Logger) *WordsStateRepo {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}

	logger := parentLogger.With().Str("component", "word_state_repo").Logger()

	return &WordsStateRepo{
		db:     db,
		logger: &logger,
	}
}

func (r *WordsStateRepo) UpsertStatus(
	ctx context.Context,
	userID int64,
	dictWordID string,
	status domain.UserWordStatus,
) error {
	const op = "UpsertStatus"

	const query = `
		INSERT INTO user_words_state (user_id, dict_word_id, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, dict_word_id) DO UPDATE
		SET status = EXCLUDED.status;
	`

	if _, err := r.db.ExecContext(ctx, query, userID, dictWordID, string(status)); err != nil {
		return fmt.Errorf("%s failed: %w", op, err)
	}

	return nil
}
