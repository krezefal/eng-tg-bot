package telegram

import (
	"context"

	tele "gopkg.in/telebot.v4"

	"github.com/krezefal/eng-tg-bot/internal/transport/telegram/ui"
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

func (t *Server) InitRoutes(ctx context.Context, h Handlers) {
	t.bot.Handle("/start", h.Start)
	t.bot.Handle("/help", h.Help)
	t.bot.Handle(ui.MainMenuHelpText, h.Help)
	t.bot.Handle("/removeMe", h.RemoveMe)

	t.bot.Handle("/dict", h.Dict)
	t.bot.Handle(ui.MainMenuDictText, h.Dict)
	t.bot.Handle("/mydict", h.List)
	t.bot.Handle(ui.MainMenuMyDictText, h.List)

	t.bot.Handle("/sub", h.Subscribe)
	t.bot.Handle("/unsub", h.Unsubscribe)

	t.bot.Handle("/learn", h.Learn)
	t.bot.Handle(&tele.InlineButton{Unique: "learn"}, h.DecisionCallback)

	t.bot.Handle("/review", h.Review)
	t.bot.Handle(&tele.InlineButton{Unique: "rate"}, h.RateCallback)
}
