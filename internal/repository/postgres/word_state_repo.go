package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/krezefal/eng-tg-bot/internal/domain"
	"github.com/rs/zerolog"
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
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *WordsStateRepo) HasReviewWords(ctx context.Context, userID int64, dictionaryID string) (bool, error) {
	const op = "HasReviewWords"

	const query = `
		SELECT EXISTS(
			SELECT 1
			FROM user_words_state uws
			INNER JOIN dictionary_words dw ON dw.id = uws.dict_word_id
			WHERE uws.user_id = $1
				AND dw.dictionary_id = $2
				AND uws.status = 'learning'
		);
	`

	var has bool
	if err := r.db.QueryRowContext(ctx, query, userID, dictionaryID).Scan(&has); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return has, nil
}

func (r *WordsStateRepo) ListDueReviewWords(
	ctx context.Context,
	userID int64,
	dictionaryID string,
	now time.Time,
) ([]*domain.ReviewWord, error) {
	const op = "ListDueReviewWords"

	const query = `
		SELECT dw.id, dw.dictionary_id, dw.spelling, dw.transcription, dw.audio, dw.ru_translation,
		       uws.ef, uws.interval_days, uws.repetition, uws.next_review_at
		FROM user_words_state uws
		INNER JOIN dictionary_words dw ON dw.id = uws.dict_word_id
		WHERE uws.user_id = $1
			AND dw.dictionary_id = $2
			AND uws.status = 'learning'
			AND (uws.next_review_at IS NULL OR uws.next_review_at <= $3)
		ORDER BY uws.next_review_at NULLS FIRST, dw.spelling ASC;
	`

	rows, err := r.db.QueryContext(ctx, query, userID, dictionaryID, now)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	words := make([]*domain.ReviewWord, 0, 16)
	for rows.Next() {
		word, scanErr := toDomainReviewWord(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("%s: %w", op, scanErr)
		}

		words = append(words, word)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return words, nil
}

func (r *WordsStateRepo) ListAllReviewWordsByNearest(
	ctx context.Context,
	userID int64,
	dictionaryID string,
	now time.Time,
) ([]*domain.ReviewWord, error) {
	const op = "ListAllReviewWordsByNearest"

	const query = `
		SELECT dw.id, dw.dictionary_id, dw.spelling, dw.transcription, dw.audio, dw.ru_translation,
		       uws.ef, uws.interval_days, uws.repetition, uws.next_review_at
		FROM user_words_state uws
		INNER JOIN dictionary_words dw ON dw.id = uws.dict_word_id
		WHERE uws.user_id = $1
			AND dw.dictionary_id = $2
			AND uws.status = 'learning'
		ORDER BY COALESCE(uws.next_review_at, $3) ASC, dw.spelling ASC;
	`

	rows, err := r.db.QueryContext(ctx, query, userID, dictionaryID, now)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	words := make([]*domain.ReviewWord, 0, 16)
	for rows.Next() {
		word, scanErr := toDomainReviewWord(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("%s: %w", op, scanErr)
		}

		words = append(words, word)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return words, nil
}

func (r *WordsStateRepo) ApplyReviewResult(
	ctx context.Context,
	in *domain.ApplyReviewResultInput,
) error {
	const op = "ApplyReviewResult"

	if in == nil {
		return fmt.Errorf("%s: input is nil", op)
	}
	if in.Result == nil {
		return fmt.Errorf("%s: result is nil", op)
	}

	const query = `
		UPDATE user_words_state
		SET ef = $3,
			interval_days = $4,
			repetition = $5,
			last_result = $6,
			last_review_at = $7,
			next_review_at = $8
		WHERE user_id = $1
			AND dict_word_id = $2;
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		in.UserID,
		in.DictWordID,
		in.Result.EF,
		in.Result.IntervalDays,
		in.Result.Repetition,
		in.Grade,
		in.ReviewedAt,
		in.Result.NextReviewAt,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rows == 0 {
		return domain.ErrReviewNotStarted
	}

	return nil
}
