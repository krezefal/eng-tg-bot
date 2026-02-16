package usecase

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type UserRepo interface {
	CreateUser(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
}

type OnboardingUsecase struct {
	userRepo UserRepo
	logger   *zerolog.Logger
}

func NewOnboardingUsecase(userRepo UserRepo, parentLogger *zerolog.Logger) *OnboardingUsecase {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}

	logger := parentLogger.With().Str("component", "onboarding_usecase").Logger()

	return &OnboardingUsecase{
		userRepo: userRepo,
		logger:   &logger,
	}
}

func (u *OnboardingUsecase) Start(ctx context.Context, userID int64) error {
	const op = "OnboardingUsecase.Start"

	// idempotence creation
	err := u.userRepo.CreateUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s failed: %w", op, err)
	}

	u.logger.Debug().
		Int64("user_id", userID).
		Msgf("%s succeeded", op)

	return nil
}

func (u *OnboardingUsecase) RemoveMe(ctx context.Context, userID int64) error {
	const op = "OnboardingUsecase.RemoveMe"

	// idempotence deletion
	err := u.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s failed: %w", op, err)
	}

	u.logger.Debug().
		Int64("user_id", userID).
		Msgf("%s succeeded", op)

	return nil
}
