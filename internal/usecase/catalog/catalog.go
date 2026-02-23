package catalog

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

type CatalogUsecase struct {
	dictRepo DictionaryRepo
	subsRepo SubscriptionsRepo
	logger   *zerolog.Logger
}

func NewUsecase(dictRepo DictionaryRepo, subsRepo SubscriptionsRepo, parentLogger *zerolog.Logger) *CatalogUsecase {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}

	logger := parentLogger.With().Str("component", "catalog_usecase").Logger()

	return &CatalogUsecase{
		dictRepo: dictRepo,
		subsRepo: subsRepo,
		logger:   &logger,
	}
}

func (u *CatalogUsecase) PublicDictionaries(ctx context.Context) ([]domain.Dictionary, error) {
	const op = "PublicDictionaries"

	dicts, err := u.dictRepo.ListPublic(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}

	u.logger.Debug().Int("count", len(dicts)).Msgf("%s succeeded", op)

	return dicts, nil
}

func (u *CatalogUsecase) UserDictionaries(ctx context.Context, userID int64) ([]domain.Dictionary, error) {
	const op = "UserDictionaries"

	dicts, err := u.subsRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}

	u.logger.Debug().
		Int64("user_id", userID).
		Int("count", len(dicts)).
		Msgf("%s succeeded", op)

	return dicts, nil
}

func (u *CatalogUsecase) DictionaryDetails(
	ctx context.Context,
	userID int64,
	dictionaryID string,
) (*domain.DictionaryDetails, error) {
	const op = "DictionaryDetails"

	dict, err := u.dictRepo.GetByID(ctx, dictionaryID)
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}

	words, err := u.dictRepo.ListRandomPreviewWords(ctx, dictionaryID, 5)
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}

	u.logger.Debug().
		Int64("user_id", userID).
		Str("dictionary_id", dictionaryID).
		Int("words_count", len(words)).
		Msgf("%s succeeded", op)

	return &domain.DictionaryDetails{
		Dictionary: dict,
		Words:      words,
	}, nil
}
