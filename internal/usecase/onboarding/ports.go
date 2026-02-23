package onboarding

import "context"

type UserRepo interface {
	CreateUser(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
}
