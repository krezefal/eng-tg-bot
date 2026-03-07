package activedict

import "context"

type UserRepo interface {
	SetActiveDictionaryID(ctx context.Context, userID int64, dictionaryID string) error
	GetActiveDictionaryID(ctx context.Context, userID int64) (string, error)
	ClearActiveDictionaryID(ctx context.Context, userID int64) error
}
