package subscription

import "context"

type UserRepo interface {
	CreateUser(ctx context.Context, id int64) error
}

type DictionaryRepo interface {
	ExistsByID(ctx context.Context, dictionaryID string) (bool, error)
}

type SubscriptionsRepo interface {
	Subscribe(ctx context.Context, userID int64, dictionaryID string) (bool, error)
	Unsubscribe(ctx context.Context, userID int64, dictionaryID string) (bool, error)
	IsSubscribedByUser(ctx context.Context, userID int64, dictionaryID string) (bool, error)
}
