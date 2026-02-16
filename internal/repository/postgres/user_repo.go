package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"
)

type UserRepo struct {
	db     *sql.DB
	logger *zerolog.Logger
}

func NewUserRepo(db *sql.DB, parentLogger *zerolog.Logger) *UserRepo {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}

	logger := parentLogger.With().Str("component", "user_repo").Logger()

	return &UserRepo{
		db:     db,
		logger: &logger,
	}
}

func (r *UserRepo) CreateUser(ctx context.Context, id int64) error {
	const op = "UserRepo.CreateUser"

	const query = `
		INSERT INTO users (tg_id)
		VALUES ($1)
		ON CONFLICT (tg_id) DO NOTHING;
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s failed: %w", op, err)
	}

	return nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, id int64) error {
	const op = "UserRepo.DeleteUser"

	const query = `
		DELETE FROM users
		WHERE tg_id = $1;
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s failed: %w", op, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s failed: %w", op, err)
	}

	if rows == 0 {
		r.logger.Warn().
			Int64("tg_id", id).
			Msgf("%s: user not found", op)
	}

	return nil
}
