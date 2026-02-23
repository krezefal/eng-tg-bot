package subscription

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

type SubscriptionUsecase struct {
	userRepo UserRepo
	dictRepo DictionaryRepo
	subsRepo SubscriptionsRepo
	logger   *zerolog.Logger
}

func NewUsecase(
	userRepo UserRepo,
	dictRepo DictionaryRepo,
	subsRepo SubscriptionsRepo,
	parentLogger *zerolog.Logger,
) *SubscriptionUsecase {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}

	logger := parentLogger.With().Str("component", "subscription_usecase").Logger()

	return &SubscriptionUsecase{
		userRepo: userRepo,
		dictRepo: dictRepo,
		subsRepo: subsRepo,
		logger:   &logger,
	}
}

func (u *SubscriptionUsecase) Subscribe(ctx context.Context, userID int64, dictionaryID string) error {
	const op = "Subscribe"

	if err := u.ensureUserAndDictionaryByID(ctx, userID, dictionaryID); err != nil {
		return fmt.Errorf("%s failed: %w", op, err)
	}

	inserted, err := u.subsRepo.Subscribe(ctx, userID, dictionaryID)
	if err != nil {
		return fmt.Errorf("%s failed: %w", op, err)
	}
	if !inserted {
		return domain.ErrAlreadySubscribed
	}

	u.logger.Debug().
		Int64("user_id", userID).
		Str("dictionary_id", dictionaryID).
		Msgf("%s succeeded", op)

	return nil
}

func (u *SubscriptionUsecase) Unsubscribe(ctx context.Context, userID int64, dictionaryID string) error {
	const op = "Unsubscribe"

	exists, err := u.dictRepo.ExistsByID(ctx, dictionaryID)
	if err != nil {
		return fmt.Errorf("%s failed: %w", op, err)
	}
	if !exists {
		return domain.ErrDictionaryNotFound
	}

	removed, err := u.subsRepo.Unsubscribe(ctx, userID, dictionaryID)
	if err != nil {
		return fmt.Errorf("%s failed: %w", op, err)
	}
	if !removed {
		return domain.ErrSubscriptionNotFound
	}

	u.logger.Debug().
		Int64("user_id", userID).
		Str("dictionary_id", dictionaryID).
		Msgf("%s succeeded", op)

	return nil
}

func (u *SubscriptionUsecase) EnsureSubscribed(ctx context.Context, userID int64, dictionaryID string) error {
	const op = "IsSubscribed"

	exists, err := u.dictRepo.ExistsByID(ctx, dictionaryID)
	if err != nil {
		return fmt.Errorf("%s failed: %w", op, err)
	}
	if !exists {
		return domain.ErrDictionaryNotFound
	}

	subscribed, err := u.subsRepo.IsSubscribedByUser(ctx, userID, dictionaryID)
	if err != nil {
		return fmt.Errorf("%s failed: %w", op, err)
	}
	if !subscribed {
		return domain.ErrSubscriptionNotFound
	}

	u.logger.Debug().
		Int64("user_id", userID).
		Str("dictionary_id", dictionaryID).
		Msgf("%s succeeded", op)

	return nil
}

func (u *SubscriptionUsecase) ensureUserAndDictionaryByID(ctx context.Context, userID int64, dictionaryID string) error {
	// Create user for the case when the user picked removeAll before.
	// Idempotent creation.
	if err := u.userRepo.CreateUser(ctx, userID); err != nil {
		return err
	}

	exists, err := u.dictRepo.ExistsByID(ctx, dictionaryID)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrDictionaryNotFound
	}

	return nil
}
