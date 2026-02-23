package telegram

import (
	"context"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

type OnboardingUsecase interface {
	Start(ctx context.Context, userID int64) error
	RemoveMe(ctx context.Context, userID int64) error
}

type CatalogUsecase interface {
	PublicDictionaries(ctx context.Context) ([]domain.Dictionary, error)
	UserDictionaries(ctx context.Context, userID int64) ([]domain.Dictionary, error)
	DictionaryDetails(ctx context.Context, userID int64, dictionaryID string) (*domain.DictionaryDetails, error)
}

type SubscriptionUsecase interface {
	Subscribe(ctx context.Context, userID int64, dictionaryID string) error
	Unsubscribe(ctx context.Context, userID int64, dictionaryID string) error
	EnsureSubscribed(ctx context.Context, userID int64, dictionaryID string) error
}

type LearningUsecase interface {
	LearnByDictionaryNumber(ctx context.Context, userID int64, number int) (*domain.LearningWord, string, error)
	LearnByDictionaryID(ctx context.Context, userID int64, dictionaryID string) (*domain.LearningWord, error)
	AddCurrentWord(ctx context.Context, userID int64) (*domain.LearningWord, error)
	BlockCurrentWord(ctx context.Context, userID int64) (*domain.LearningWord, error)
	ActiveDictionaryID(ctx context.Context, userID int64) (string, error)
	Back(ctx context.Context, userID int64) error
}

type ReviewUsecase interface {
	Review(ctx context.Context, userID int64, dictionaryID string) error
	RateCallback(ctx context.Context, userID int64, wordID string, rate int) error
}
