package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

type SubscriptionsRepo struct {
	db     *sql.DB
	logger *zerolog.Logger
}

func NewSubscriptionsRepo(db *sql.DB, parentLogger *zerolog.Logger) *SubscriptionsRepo {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}

	logger := parentLogger.With().Str("component", "subscriptions_repo").Logger()

	return &SubscriptionsRepo{
		db:     db,
		logger: &logger,
	}
}

func (r *SubscriptionsRepo) Subscribe(ctx context.Context, userID int64, dictionaryID string) (bool, error) {
	const op = "Subscribe"

	const query = `
		INSERT INTO user_dictionaries (user_id, dictionary_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, dictionary_id) DO NOTHING;
	`

	res, err := r.db.ExecContext(ctx, query, userID, dictionaryID)
	if err != nil {
		return false, fmt.Errorf("%s failed: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s failed: %w", op, err)
	}

	return rows > 0, nil
}

func (r *SubscriptionsRepo) Unsubscribe(ctx context.Context, userID int64, dictionaryID string) (bool, error) {
	const op = "Unsubscribe"

	const query = `
		DELETE FROM user_dictionaries
		WHERE user_id = $1 AND dictionary_id = $2;
	`

	res, err := r.db.ExecContext(ctx, query, userID, dictionaryID)
	if err != nil {
		return false, fmt.Errorf("%s failed: %w", op, err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s failed: %w", op, err)
	}

	return rows > 0, nil
}

func (r *SubscriptionsRepo) ListByUser(ctx context.Context, userID int64) ([]domain.Dictionary, error) {
	const op = "ListByUser"

	const query = `
		SELECT d.id, d.title, d.description, d.mode, d.author, d.created_at
		FROM user_dictionaries ud
		INNER JOIN dictionaries d ON d.id = ud.dictionary_id
		WHERE ud.user_id = $1
		ORDER BY ud.subscribed_at DESC, d.title ASC;
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}
	defer rows.Close()

	dictionaries := make([]domain.Dictionary, 0, 16)
	for rows.Next() {
		d, scanErr := toDomainDictionary(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("%s failed: %w", op, scanErr)
		}

		dictionaries = append(dictionaries, *d)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}

	return dictionaries, nil
}

func (r *SubscriptionsRepo) IsSubscribedByUser(ctx context.Context, userID int64, dictionaryID string) (bool, error) {
	const op = "IsSubscribedByUser"

	const query = `
		SELECT EXISTS(
		SELECT 1
		FROM user_dictionaries
		WHERE user_id = $1 AND dictionary_id = $2
		);
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, dictionaryID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s failed: %w", op, err)
	}

	return exists, nil
}
