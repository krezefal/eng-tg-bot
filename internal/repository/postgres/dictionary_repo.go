package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

type DictionaryRepo struct {
	db     *sql.DB
	logger *zerolog.Logger
}

func NewDictionaryRepo(db *sql.DB, parentLogger *zerolog.Logger) *DictionaryRepo {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}

	logger := parentLogger.With().Str("component", "dictionary_repo").Logger()

	return &DictionaryRepo{
		db:     db,
		logger: &logger,
	}
}

func (r *DictionaryRepo) ListPublic(ctx context.Context) ([]domain.Dictionary, error) {
	const op = "ListPublic"

	const query = `
		SELECT id, title, description, mode, author, created_at
		FROM dictionaries
		ORDER BY created_at DESC, title ASC;
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}
	defer rows.Close()

	dictionaries := make([]domain.Dictionary, 0, 32)
	for rows.Next() {
		d, scanErr := toDomainDictionary(rows)
		if scanErr != nil {
			err = scanErr
			return nil, fmt.Errorf("%s failed: %w", op, err)
		}

		dictionaries = append(dictionaries, *d)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}

	return dictionaries, nil
}

func (r *DictionaryRepo) GetByID(ctx context.Context, dictionaryID string) (*domain.Dictionary, error) {
	const op = "GetByID"

	const query = `
		SELECT id, title, description, mode, author, created_at
		FROM dictionaries
		WHERE id = $1;
	`

	row := r.db.QueryRowContext(ctx, query, dictionaryID)
	dict, err := toDomainDictionary(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDictionaryNotFound
		}

		return nil, fmt.Errorf("%s failed: %w", op, err)
	}

	return dict, nil
}

func (r *DictionaryRepo) ExistsByID(ctx context.Context, dictionaryID string) (bool, error) {
	const op = "ExistsByID"

	const query = `
		SELECT EXISTS(
			SELECT 1
			FROM dictionaries
			WHERE id = $1
		);
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, dictionaryID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s failed: %w", op, err)
	}

	return exists, nil
}

func (r *DictionaryRepo) ListRandomPreviewWords(
	ctx context.Context,
	dictionaryID string,
	limit int,
) ([]domain.DictionaryWordPreview, error) {
	const op = "ListRandomPreviewWords"

	const query = `
		SELECT spelling, ru_translation
		FROM dictionary_words
		WHERE dictionary_id = $1
		ORDER BY random()
		LIMIT $2;
	`

	rows, err := r.db.QueryContext(ctx, query, dictionaryID, limit)
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}
	defer rows.Close()

	words := make([]domain.DictionaryWordPreview, 0, limit)
	for rows.Next() {
		w, scanErr := toDomainDictionaryWordPreview(rows)
		if scanErr != nil {
			err = scanErr
			return nil, fmt.Errorf("%s failed: %w", op, err)
		}

		words = append(words, *w)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}

	return words, nil
}

// TODO: good place for caching batch of untracked words not to pick from DB
// every time.
func (r *DictionaryRepo) PickRandomUntrackedWord(
	ctx context.Context,
	userID int64,
	dictionaryID string,
) (*domain.LearningWord, error) {
	const op = "PickRandomUntrackedWord"

	const query = `
		SELECT dw.id, dw.dictionary_id, dw.spelling, dw.transcription, dw.audio, dw.ru_translation
		FROM dictionary_words dw
		LEFT JOIN user_words_state uws
			ON uws.dict_word_id = dw.id AND uws.user_id = $1
		WHERE dw.dictionary_id = $2
			AND uws.dict_word_id IS NULL
		ORDER BY random()
		LIMIT 1;
	`

	row := r.db.QueryRowContext(ctx, query, userID, dictionaryID)
	word, err := toDomainLearningWord(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("%s failed: %w", op, err)
	}

	return word, nil
}
