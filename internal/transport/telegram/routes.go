package telegram

import (
	"context"

	tele "gopkg.in/telebot.v4"

	"github.com/krezefal/eng-tg-bot/internal/transport/telegram/ui"
)

type Handlers interface {
	// Onboarding
	Start(c tele.Context) error
	Help(c tele.Context) error
	RemoveMe(c tele.Context) error

	// Catalog
	Dict(c tele.Context) error
	MyDict(c tele.Context) error
	DictDetails(c tele.Context) error

	// Subscription
	Subscribe(c tele.Context) error
	Unsubscribe(c tele.Context) error
	ConfirmUnsubscribe(c tele.Context) error
	RejectUnsubscribe(c tele.Context) error

	// Learning
	LearnByDictNum(c tele.Context) error
	LearnByDictID(c tele.Context) error
	LearningAction(c tele.Context) error

	// Review
	Review(c tele.Context) error
	RateCallback(c tele.Context) error
}

func (t *Server) InitRoutes(_ context.Context, h Handlers) {
	// Onboarding
	t.bot.Handle("/start", h.Start)
	t.bot.Handle("/help", h.Help)
	t.bot.Handle(ui.MainMenuHelpText, h.Help)
	t.bot.Handle("/removeMe", h.RemoveMe)

	// Catalog
	t.bot.Handle("/dict", h.Dict)
	t.bot.Handle(ui.MainMenuDictText, h.Dict)
	t.bot.Handle(&tele.InlineButton{Unique: "to_dicts"}, h.Dict)
	t.bot.Handle("/mydict", h.MyDict)
	t.bot.Handle(ui.MainMenuMyDictText, h.MyDict)
	t.bot.Handle(&tele.InlineButton{Unique: "dict_details"}, h.DictDetails)

	// Subscription
	t.bot.Handle(&tele.InlineButton{Unique: "dict_subscribe"}, h.Subscribe)
	t.bot.Handle(&tele.InlineButton{Unique: "dict_unsubscribe"}, h.Unsubscribe)
	t.bot.Handle(&tele.InlineButton{Unique: "dict_confirm_unsubscribe"}, h.ConfirmUnsubscribe)
	t.bot.Handle(&tele.InlineButton{Unique: "dict_reject_unsubscribe"}, h.RejectUnsubscribe)

	// Learning
	t.bot.Handle("/learn", h.LearnByDictNum)
	t.bot.Handle(&tele.InlineButton{Unique: "dict_learn"}, h.LearnByDictID)
	t.bot.Handle(ui.LearnAddText, h.LearningAction)
	t.bot.Handle(ui.LearnBlockText, h.LearningAction)
	t.bot.Handle(ui.LearnReviewText, h.LearningAction)
	t.bot.Handle(ui.LearnBackText, h.LearningAction)

	// Review
	t.bot.Handle("/review", h.Review)
	t.bot.Handle(&tele.InlineButton{Unique: "dict_review"}, h.Review)
	t.bot.Handle(&tele.InlineButton{Unique: "rate"}, h.RateCallback)
}
