package telebot

import (
	"context"

	tele "gopkg.in/telebot.v4"
)

type Handlers interface {
	Start(c tele.Context) error
	Help(c tele.Context) error
	RemoveMe(c tele.Context) error

	Dict(c tele.Context) error
	List(c tele.Context) error

	Subscribe(c tele.Context) error
	Unsubscribe(c tele.Context) error

	Learn(c tele.Context) error
	DecisionCallback(c tele.Context) error

	Review(c tele.Context) error
	RateCallback(c tele.Context) error
}

func (t *TelebotServer) InitRoutes(ctx context.Context, h Handlers) {
	t.telebot.Handle("/start", h.Start)
	t.telebot.Handle("/help", h.Help)
	t.telebot.Handle("/помощь", h.Help)
	t.telebot.Handle("/removeMe", h.RemoveMe)

	t.telebot.Handle("/словари", h.Dict)
	t.telebot.Handle("/моисловари", h.List)

	t.telebot.Handle("/подписаться", h.Subscribe)
	t.telebot.Handle("/отписаться", h.Unsubscribe)

	t.telebot.Handle("/учить", h.Learn)
	t.telebot.Handle(&tele.InlineButton{Unique: "learn"}, h.DecisionCallback)

	t.telebot.Handle("/повторить", h.Review)
	t.telebot.Handle(&tele.InlineButton{Unique: "rate"}, h.RateCallback)
}
