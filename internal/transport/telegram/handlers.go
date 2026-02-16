package telegram

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"

	"github.com/krezefal/eng-tg-bot/internal/transport/telegram/ui"
)

const handlerCtxTimeout = 5 * time.Second

var _ Handlers = (*BotHandlers)(nil)

// TODO: add recovering from panics somewhere
type BotHandlers struct {
	onboardUC OnboardingUsecase
	catalogUC CatalogLUsecase
	subsUC    SubscriptionUsecase
	learnUC   LearningUsecase
	reviewUC  ReviewUsecase
	logger    *zerolog.Logger
}

func NewHandler(
	onboardUC OnboardingUsecase,
	catalogUC CatalogLUsecase,
	subsUC SubscriptionUsecase,
	learnUC LearningUsecase,
	reviewUC ReviewUsecase,
	parentLogger *zerolog.Logger,
) *BotHandlers {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}

	// TODO: uncomment
	//if onboardUC == nil {
	//	panic("OnboardingUsecase cannot be nil")
	//}
	//if catalogUC == nil {
	//	panic("CatalogLUsecase cannot be nil")
	//}
	//if subsUC == nil {
	//	panic("SubscriptionUsecase cannot be nil")
	//}
	//if learnUC == nil {
	//	panic("LearningUsecase cannot be nil")
	//}
	//if reviewUC == nil {
	//	panic("ReviewUsecase cannot be nil")
	//}

	logger := parentLogger.With().Str("component", "telegram_handler").Logger()

	return &BotHandlers{
		onboardUC: onboardUC,
		catalogUC: catalogUC,
		subsUC:    subsUC,
		learnUC:   learnUC,
		reviewUC:  reviewUC,
		logger:    &logger,
	}
}

func (h *BotHandlers) Start(c tele.Context) error {
	const op = "Start"

	ctx, cancel := context.WithTimeout(context.Background(), handlerCtxTimeout)
	defer cancel()

	userID := c.Sender().ID
	updateID := c.Update().ID
	username := c.Sender().Username
	messageID := 0
	if msg := c.Message(); msg != nil {
		messageID = msg.ID
	}

	// TODO: remove personal data from logs after alfa-test
	h.logger.Debug().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Int("message_id", messageID).
		Msgf("handling %s", op)

	if err := h.onboardUC.Start(ctx, userID); err != nil {
		h.logger.Error().
			Err(err).
			Int64("user_id", userID).
			Msgf("%s failed", op)

		return c.Send(ui.InternalErrorMessage)
	}

	return c.Send(ui.WelcomeMessage, ui.BuildMainMenuKeyboard())
}

func (h *BotHandlers) Help(c tele.Context) error {
	return c.Send(ui.HelpMessage, ui.BuildMainMenuKeyboard())
}

func (h *BotHandlers) RemoveMe(c tele.Context) error {
	const op = "RemoveMe"

	ctx, cancel := context.WithTimeout(context.Background(), handlerCtxTimeout)
	defer cancel()

	userID := c.Sender().ID
	updateID := c.Update().ID
	username := c.Sender().Username
	messageID := 0
	if msg := c.Message(); msg != nil {
		messageID = msg.ID
	}

	// TODO: remove personal data from logs after alfa-test
	h.logger.Debug().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Int("message_id", messageID).
		Msgf("handling %s", op)

	if err := h.onboardUC.RemoveMe(ctx, userID); err != nil {
		h.logger.Error().
			Err(err).
			Int64("user_id", userID).
			Msgf("%s failed", op)

		return c.Send(ui.InternalErrorMessage)
	}

	return c.Send(ui.RemoveMessage)
}

func (h *BotHandlers) Dict(c tele.Context) error {
	return c.Send("WIP")
}

func (h *BotHandlers) List(c tele.Context) error {
	return c.Send("WIP")
}

func (h *BotHandlers) Subscribe(c tele.Context) error {
	return c.Send("WIP")
}

func (h *BotHandlers) Unsubscribe(c tele.Context) error {
	return c.Send("WIP")
}

func (h *BotHandlers) Learn(c tele.Context) error {
	return c.Send("WIP")
}

func (h *BotHandlers) DecisionCallback(c tele.Context) error {
	return c.Send("WIP")
}

func (h *BotHandlers) Review(c tele.Context) error {
	return c.Send("WIP")

	//args := c.Args()
	//if len(args) == 0 {
	//	return c.Send("Usage: /repeat <vocabulary>")
	//}
	//return c.Send(strings.ToUpper(args[0]))
}

func (h *BotHandlers) RateCallback(c tele.Context) error {
	return c.Send("WIP")
}

//func (h *BotHandlers) Review(c tele.Context) error {
//	userID := c.Sender().ID
//	dictID := parseDictID(c)
//
//	card, err := h.reviewLogic.Next(c, userID, dictID)
//	if err != nil {
//		return c.Send("Не удалось подобрать слово")
//	}
//
//	kb := ui.BuildRateKeyboard(card.ID)
//
//	return c.Send(formatCard(card), kb)
//}
//
//func (h *BotHandlers) RateCallback(c tele.Context) error {
//	userID := c.Sender().ID
//
//	// data: "<wordID>:<grade>"
//	data := c.Data()
//
//	parts := strings.Split(data, ":")
//	if len(parts) != 2 {
//		_ = c.Respond() // убрать "часики"
//		return nil
//	}
//
//	wordID := parts[0]
//
//	grade, err := strconv.Atoi(parts[1])
//	if err != nil {
//		_ = c.Respond()
//		return nil
//	}
//
//	// вызываем usecase (без tele.Context)
//	err = h.reviewLogic.Rate(c, userID, wordID, grade)
//	if err != nil {
//		_ = c.Respond()
//		return c.Send("Ошибка при сохранении оценки")
//	}
//
//	// убрать "часики" у кнопки
//	_ = c.Respond()
//
//	return c.Send("Оценка сохранена ✅")
//}
