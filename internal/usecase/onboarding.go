package logic

import (
	"context"

	"github.com/rs/zerolog"
)

type UserRepo interface {
	CreateUser(ctx context.Context, id int64) error
}

type OnboardingLogic struct {
	userRepo UserRepo
	logger   *zerolog.Logger
}

func NewOnboardingLogic() {

}
