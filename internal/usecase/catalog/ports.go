package catalog

import (
	"context"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

type DictionaryRepo interface {
	ListPublic(ctx context.Context) ([]domain.Dictionary, error)
	GetByID(ctx context.Context, dictionaryID string) (*domain.Dictionary, error)
	ListRandomPreviewWords(ctx context.Context, dictionaryID string, limit int) ([]domain.DictionaryWordPreview, error)
}

type SubscriptionsRepo interface {
	ListByUser(ctx context.Context, userID int64) ([]domain.Dictionary, error)
}
