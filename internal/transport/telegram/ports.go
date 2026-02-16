package telegram

import (
	"context"

	tele "gopkg.in/telebot.v4"
)

type OnboardingUsecase interface {
	Start(ctx context.Context, userID int64) error
	RemoveMe(ctx context.Context, userID int64) error
}

type CatalogLUsecase interface {
	Dict(c tele.Context) error
	List(c tele.Context) error
}

type SubscriptionUsecase interface {
	Subscribe(c tele.Context) error
	Unsubscribe(c tele.Context) error
}

type LearningUsecase interface {
	Learn(c tele.Context) error
	DecisionCallback(c tele.Context) error
}

type ReviewUsecase interface {
	Review(c tele.Context) error
	RateCallback(c tele.Context) error
}
