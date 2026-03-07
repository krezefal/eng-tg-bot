package learning

import (
	"context"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

type DictionaryRepo interface {
	ExistsByID(ctx context.Context, dictionaryID string) (bool, error)
	PickRandomUntrackedWord(ctx context.Context, userID int64, dictionaryID string) (*domain.LearningWord, error)
}

type SubscriptionsRepo interface {
	ListByUser(ctx context.Context, userID int64) ([]domain.Dictionary, error)
	IsSubscribedByUser(ctx context.Context, userID int64, dictionaryID string) (bool, error)
	MarkLearningStarted(ctx context.Context, userID int64, dictionaryID string) error
}

type WordStateRepo interface {
	UpsertStatus(ctx context.Context, userID int64, dictWordID string, status domain.UserWordStatus) error
}
