package activedict

import (
	"context"
	"fmt"

	_ "github.com/krezefal/eng-tg-bot/internal/domain"
	"github.com/rs/zerolog"
)

type ActiveDictUseCase struct {
	userRepo UserRepo
	logger   *zerolog.Logger
}

func NewActiveDictUseCase(
	userRepo UserRepo,
	parentLogger *zerolog.Logger,
) *ActiveDictUseCase {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}

	logger := parentLogger.With().Str("component", "learning_usecase").Logger()

	return &ActiveDictUseCase{
		userRepo: userRepo,
		logger:   &logger,
	}
}

func (u *ActiveDictUseCase) GetActiveDictionaryID(ctx context.Context, userID int64) (string, error) {
	const op = "ActiveDictionaryID"

	dictionaryID, err := u.userRepo.GetActiveDictionaryID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return dictionaryID, nil
}

func (u *ActiveDictUseCase) SetActiveDictionaryID(ctx context.Context, userID int64, dictionaryID string) error {
	const op = "ActiveDictionaryID"

	if err := u.userRepo.SetActiveDictionaryID(ctx, userID, dictionaryID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *ActiveDictUseCase) ClearActiveDictionaryID(ctx context.Context, userID int64) error {
	const op = "ActiveDictionaryID"

	if err := u.userRepo.ClearActiveDictionaryID(ctx, userID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
