package telebot

import (
	"strconv"
	"strings"

	"github.com/krezefal/eng-tg-bot/internal/models"
	telebotServer "github.com/krezefal/eng-tg-bot/internal/server/telebot"
	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"
)

var _ telebotServer.Handlers = (*Handlers)(nil)

type Handlers struct {
	onboardLogic models.OnboardingLogic
	catalogLogic models.CatalogLogic
	subsLogic    models.SubscriptionLogic
	learnLogic   models.LearningLogic
	reviewLogic  models.ReviewLogic
	logger       *zerolog.Logger
}

func NewHandler(
	onboardLogic models.OnboardingLogic,
	catalogLogic models.CatalogLogic,
	subsLogic models.SubscriptionLogic,
	learnLogic models.LearningLogic,
	reviewLogic models.ReviewLogic,
	logger *zerolog.Logger,
) *Handlers {
	if logger == nil {
		panic("logger cannot be nil")
	}

	if onboardLogic == nil {
		panic("onboardLogic cannot be nil")
	}
	if catalogLogic == nil {
		panic("catalogLogic cannot be nil")
	}
	if subsLogic == nil {
		panic("subsLogic cannot be nil")
	}
	if learnLogic == nil {
		panic("learnLogic cannot be nil")
	}
	if reviewLogic == nil {
		panic("reviewLogic cannot be nil")
	}

	return &Handlers{
		onboardLogic: onboardLogic,
		catalogLogic: catalogLogic,
		subsLogic:    subsLogic,
		learnLogic:   learnLogic,
		reviewLogic:  reviewLogic,
		logger:       logger,
	}
}

func (h *Handlers) Start(c tele.Context) error {
	return h.onboardLogic.Start(c)
}

func (h *Handlers) Help(c tele.Context) error {
	return h.onboardLogic.Help(c)
}

func (h *Handlers) RemoveMe(c tele.Context) error {
	return h.logic.RemoveMe(c)
}

func (h *Handlers) Dict(c tele.Context) error {
	return h.logic.Dict(c)
}

func (h *Handlers) List(c tele.Context) error {
	return h.logic.List(c)
}

func (h *Handlers) Subscribe(c tele.Context) error {
	return h.logic.Subscribe(c)
}

func (h *Handlers) Unsubscribe(c tele.Context) error {
	return h.logic.Unsubscribe(c)
}

func (h *Handlers) Learn(c tele.Context) error {
	return h.logic.Learn(c)
}

func (h *Handlers) DecisionCallback(c tele.Context) error {
	return h.logic.DecisionCallback(c)
}

func (h *Handlers) Review(c tele.Context) error {
	userID := c.Sender().ID
	dictID := parseDictID(c)

	card, err := h.reviewLogic.Next(c, userID, dictID)
	if err != nil {
		return c.Send("Не удалось подобрать слово")
	}

	kb := models.BuildRateKeyboard(card.ID)

	return c.Send(formatCard(card), kb)
}

func (h *Handlers) RateCallback(c tele.Context) error {
	userID := c.Sender().ID

	// data: "<wordID>:<grade>"
	data := c.Data()

	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		_ = c.Respond() // убрать "часики"
		return nil
	}

	wordID := parts[0]

	grade, err := strconv.Atoi(parts[1])
	if err != nil {
		_ = c.Respond()
		return nil
	}

	// вызываем usecase (без tele.Context)
	err = h.reviewLogic.Rate(c, userID, wordID, grade)
	if err != nil {
		_ = c.Respond()
		return c.Send("Ошибка при сохранении оценки")
	}

	// убрать "часики" у кнопки
	_ = c.Respond()

	return c.Send("Оценка сохранена ✅")
}
