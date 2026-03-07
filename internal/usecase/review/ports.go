package review

import (
	"context"
	"time"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

type DictionaryRepo interface {
	ExistsByID(ctx context.Context, dictionaryID string) (bool, error)
}

type SubscriptionsRepo interface {
	ListByUser(ctx context.Context, userID int64) ([]domain.Dictionary, error)
	IsSubscribedByUser(ctx context.Context, userID int64, dictionaryID string) (bool, error)
}

type WordsStateRepo interface {
	HasReviewWords(ctx context.Context, userID int64, dictionaryID string) (bool, error)
	ListDueReviewWords(ctx context.Context, userID int64, dictionaryID string, now time.Time) ([]*domain.ReviewWord, error)
	ListAllReviewWordsByNearest(ctx context.Context, userID int64, dictionaryID string, now time.Time) ([]*domain.ReviewWord, error)
	ApplyReviewResult(ctx context.Context, in *domain.ApplyReviewResultInput) error
}
