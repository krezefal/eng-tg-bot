package telegram

import (
	"context"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

type OnboardingUsecase interface {
	Start(ctx context.Context, userID int64, username string) error
	RemoveMe(ctx context.Context, userID int64) error
}

type CatalogUsecase interface {
	PublicDictionaries(ctx context.Context) ([]domain.Dictionary, error)
	UserDictionaries(ctx context.Context, userID int64) ([]domain.Dictionary, error)
	DictionaryDetails(ctx context.Context, userID int64, dictionaryID string) (*domain.DictionaryDetails, error)
}

type SubscriptionUsecase interface {
	Subscribe(ctx context.Context, userID int64, username, dictionaryID string) error
	Unsubscribe(ctx context.Context, userID int64, dictionaryID string) error
	EnsureSubscribed(ctx context.Context, userID int64, dictionaryID string) error
}

type LearningUsecase interface {
	LearnByDictionaryNumber(ctx context.Context, userID int64, number int) (*domain.LearningWord, string, error)
	LearnByDictionaryID(ctx context.Context, userID int64, dictionaryID string) (*domain.LearningWord, error)
	AddCurrentWord(ctx context.Context, userID int64) (*domain.LearningWord, error)
	BlockCurrentWord(ctx context.Context, userID int64) (*domain.LearningWord, error)
	Back(ctx context.Context, userID int64) error
}

type ReviewUsecase interface {
	PrepareByDictionaryNumber(ctx context.Context, userID int64, number int) (string, error)
	PrepareByDictionaryID(ctx context.Context, userID int64, dictionaryID string) error
	StartDueRound(ctx context.Context, userID int64, dictionaryID string) (*domain.ReviewWord , error)
	StartForceRound(ctx context.Context, userID int64, dictionaryID string) (*domain.ReviewWord, error)
	RateCurrent(ctx context.Context, userID int64, grade int) (*domain.ReviewWord, string, error)
	Stop(ctx context.Context, userID int64) error
}

// TODO: move ActiveDictionaryID from 2 usecases above to this one. [x]
type ActiveDictionaryUsecase interface {
	GetActiveDictionaryID(ctx context.Context, userID int64) (string, error)
	SetActiveDictionaryID(ctx context.Context, userID int64, dictionaryID string) error
	ClearActiveDictionaryID(ctx context.Context, userID int64) error
}
