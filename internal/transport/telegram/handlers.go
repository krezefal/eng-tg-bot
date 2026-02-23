package telegram

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/krezefal/eng-tg-bot/internal/transport/telegram/mapper"
	"github.com/rs/zerolog"
	tele "gopkg.in/telebot.v4"

	"github.com/krezefal/eng-tg-bot/internal/domain"
	"github.com/krezefal/eng-tg-bot/internal/transport/telegram/ui"
)

const handlerCtxTimeout = 5 * time.Second

var _ Handlers = (*BotHandlers)(nil)

// TODO: add recovering from panics somewhere
type BotHandlers struct {
	onboardUC OnboardingUsecase
	catalogUC CatalogUsecase
	subsUC    SubscriptionUsecase
	learnUC   LearningUsecase
	reviewUC  ReviewUsecase
	logger    *zerolog.Logger
}

func NewHandler(
	onboardUC OnboardingUsecase,
	catalogUC CatalogUsecase,
	subsUC SubscriptionUsecase,
	learnUC LearningUsecase,
	reviewUC ReviewUsecase,
	parentLogger *zerolog.Logger,
) *BotHandlers {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}

	if onboardUC == nil {
		panic("OnboardingUsecase cannot be nil")
	}
	if catalogUC == nil {
		panic("CatalogUsecase cannot be nil")
	}
	if subsUC == nil {
		panic("SubscriptionUsecase cannot be nil")
	}
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

	// TODO: remove personal data from logs after alfa-test
	ctxLogger := h.logger.With().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Logger()

	ctxLogger.Debug().Msgf("handling %s", op)

	if err := h.onboardUC.Start(ctx, userID); err != nil {
		ctxLogger.Error().Err(err).Msgf("%s failed", op)

		return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
	}

	ctxLogger.Debug().Msgf("%s handled", op)

	return c.Send(ui.WelcomeMsg, ui.BuildMainMenuReplyKb())
}

func (h *BotHandlers) Help(c tele.Context) error {
	return c.Send(ui.HelpMsg, ui.BuildMainMenuReplyKb())
}

func (h *BotHandlers) RemoveMe(c tele.Context) error {
	const op = "RemoveMe"

	ctx, cancel := context.WithTimeout(context.Background(), handlerCtxTimeout)
	defer cancel()

	userID := c.Sender().ID
	updateID := c.Update().ID
	username := c.Sender().Username

	// TODO: remove personal data from logs after alfa-test
	ctxLogger := h.logger.With().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Logger()

	ctxLogger.Debug().Msgf("handling %s", op)

	if err := h.onboardUC.RemoveMe(ctx, userID); err != nil {
		ctxLogger.Error().Err(err).Msgf("%s failed", op)

		return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
	}

	ctxLogger.Info().Msgf("%s handled", op)

	return c.Send(ui.RemoveMsg, ui.BuildMainMenuReplyKb())
}

func (h *BotHandlers) Dict(c tele.Context) error {
	const op = "Dict"

	ctx, cancel := context.WithTimeout(context.Background(), handlerCtxTimeout)
	defer cancel()
	if c.Callback() != nil {
		defer func() {
			_ = c.Respond()
		}()
	}

	userID := c.Sender().ID
	updateID := c.Update().ID
	username := c.Sender().Username

	// TODO: remove personal data from logs after alfa-test
	ctxLogger := h.logger.With().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Logger()

	ctxLogger.Debug().Msgf("handling %s", op)

	dictionaries, err := h.catalogUC.PublicDictionaries(ctx)
	if err != nil {
		ctxLogger.Error().Err(err).Msgf("%s failed", op)

		return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
	}

	if len(dictionaries) == 0 {
		ctxLogger.Warn().Msgf("%s: no public dicts found", op)

		return c.Send(ui.PublicDictionariesEmptyMsg, ui.BuildMainMenuReplyKb())
	}
	if err = c.Send(ui.PublicDictionariesHeaderMsg, ui.BuildMainMenuReplyKb()); err != nil {
		ctxLogger.Error().Err(err).Msgf("%s failed sent main_menu_kb", op)

		return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
	}
	for _, d := range dictionaries {
		if err = c.Send(
			ui.FormatDictionaryCard(d),
			&tele.SendOptions{
				ParseMode:   tele.ModeHTML,
				ReplyMarkup: ui.BuildPublicDictionaryInlineKb(d.ID),
			},
		); err != nil {
			ctxLogger.Error().Err(err).Msgf("%s failed send dict", op)

			return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
		}
	}

	ctxLogger.Debug().Msgf("%s handled", op)

	return nil
}

func (h *BotHandlers) MyDict(c tele.Context) error {
	const op = "MyDict"

	ctx, cancel := context.WithTimeout(context.Background(), handlerCtxTimeout)
	defer cancel()
	if c.Callback() != nil {
		defer func() {
			_ = c.Respond()
		}()
	}

	userID := c.Sender().ID
	updateID := c.Update().ID
	username := c.Sender().Username

	// TODO: remove personal data from logs after alfa-test
	ctxLogger := h.logger.With().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Logger()

	ctxLogger.Debug().Msgf("handling %s", op)

	dictionaries, err := h.catalogUC.UserDictionaries(ctx, userID)
	if err != nil {
		ctxLogger.Error().Err(err).Msgf("%s failed", op)

		return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
	}

	if len(dictionaries) == 0 {
		ctxLogger.Debug().Msgf("%s: user doesn't have subscribed dicts", op)

		return c.Send(ui.UserDictionariesEmptyMsg, ui.BuildMainMenuReplyKb())
	}
	if err = c.Send(ui.UserDictionariesHeaderMsg, ui.BuildMainMenuReplyKb()); err != nil {
		ctxLogger.Error().Err(err).Msgf("%s failed", op)

		return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
	}
	for i, d := range dictionaries {
		if err = c.Send(
			ui.FormatSubscribedDictionaryCard(i+1, d),
			&tele.SendOptions{
				ParseMode:   tele.ModeHTML,
				ReplyMarkup: ui.BuildUserDictionaryInlineKb(d.ID),
			},
		); err != nil {
			ctxLogger.Error().Err(err).Msgf("%s failed", op)

			return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
		}
	}

	ctxLogger.Debug().Msgf("%s handled", op)

	return nil
}

func (h *BotHandlers) DictDetails(c tele.Context) error {
	const op = "DictDetails"

	ctx, cancel := context.WithTimeout(context.Background(), handlerCtxTimeout)
	defer cancel()
	if c.Callback() != nil {
		defer func() {
			_ = c.Respond()
		}()
	}

	userID := c.Sender().ID
	updateID := c.Update().ID
	username := c.Sender().Username

	ctxLogger := h.logger.With().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Logger()

	ctxLogger.Debug().Msgf("handling %s", op)

	dictionaryID := extractCallbackDictionaryID(c)
	if dictionaryID == "" {
		// TODO: alert here
		ctxLogger.Error().Msgf("%s: dictionary_id is empty", op)

		return c.Send(ui.DictionaryNotFoundMsg, ui.BuildMainMenuReplyKb())
	}

	ctxLogger = ctxLogger.With().Str("dictionary_id", dictionaryID).Logger()

	details, err := h.catalogUC.DictionaryDetails(ctx, userID, dictionaryID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDictionaryNotFound):
			// TODO: alert here
			ctxLogger.Error().Err(err).Msgf("%s: unable to find dictionary by dictionary_id", op)

			return c.Send(ui.DictionaryNotFoundMsg, ui.BuildMainMenuReplyKb())

		default:
			ctxLogger.Error().Err(err).Msgf("%s failed", op)

			return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
		}
	}

	ctxLogger.Debug().Msgf("%s handled", op)

	return c.Send(
		ui.FormatDictionaryDetails(*details.Dictionary, details.Words),
		&tele.SendOptions{
			ParseMode:   tele.ModeHTML,
			ReplyMarkup: ui.BuildDictionaryDetailsInlineKb(details.Dictionary.ID),
		},
	)
}

func (h *BotHandlers) Subscribe(c tele.Context) error {
	const op = "Subscribe"

	ctx, cancel := context.WithTimeout(context.Background(), handlerCtxTimeout)
	defer cancel()
	if c.Callback() != nil {
		defer func() {
			_ = c.Respond()
		}()
	}

	userID := c.Sender().ID
	updateID := c.Update().ID
	username := c.Sender().Username

	// TODO: remove personal data from logs after alfa-test
	ctxLogger := h.logger.With().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Logger()

	ctxLogger.Debug().Msgf("handling %s", op)

	dictionaryID := extractCallbackDictionaryID(c)
	if dictionaryID == "" {
		// TODO: alert here
		ctxLogger.Error().Msgf("%s: dictionary_id is empty", op)

		return c.Send(ui.DictionaryNotFoundMsg, ui.BuildMainMenuReplyKb())
	}

	ctxLogger = ctxLogger.With().Str("dictionary_id", dictionaryID).Logger()

	if err := h.subsUC.Subscribe(ctx, userID, dictionaryID); err != nil {
		switch {
		case errors.Is(err, domain.ErrDictionaryNotFound):
			// TODO: alert here
			ctxLogger.Error().Err(err).Msgf("%s: unable to find dictionary by dictionary_id", op)

			return c.Send(ui.DictionaryNotFoundMsg, ui.BuildMainMenuReplyKb())

		case errors.Is(err, domain.ErrAlreadySubscribed):
			ctxLogger.Debug().Msgf("%s: already subscribed", op)

			return c.Send(ui.DictionaryAlreadySubscribedMsg, ui.BuildMainMenuReplyKb())

		default:
			ctxLogger.Error().Err(err).Msgf("%s failed", op)

			return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
		}
	}

	ctxLogger.Debug().Msgf("%s handled", op)

	return c.Send(ui.DictionarySubscribedMsg, ui.BuildMainMenuReplyKb())
}

func (h *BotHandlers) Unsubscribe(c tele.Context) error {
	const op = "Unsubscribe"

	ctx, cancel := context.WithTimeout(context.Background(), handlerCtxTimeout)
	defer cancel()
	if c.Callback() != nil {
		defer func() {
			_ = c.Respond()
		}()
	}

	userID := c.Sender().ID
	updateID := c.Update().ID
	username := c.Sender().Username

	// TODO: remove personal data from logs after alfa-test
	ctxLogger := h.logger.With().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Logger()

	ctxLogger.Debug().Msgf("handling %s", op)

	dictionaryID := extractCallbackDictionaryID(c)
	if dictionaryID == "" {
		// TODO: alert here
		ctxLogger.Error().Msgf("%s: dictionary_id is empty", op)

		return c.Send(ui.DictionaryNotFoundMsg, ui.BuildMainMenuReplyKb())
	}

	ctxLogger = ctxLogger.With().Str("dictionary_id", dictionaryID).Logger()

	if err := h.subsUC.EnsureSubscribed(ctx, userID, dictionaryID); err != nil {
		switch {
		case errors.Is(err, domain.ErrDictionaryNotFound):
			// TODO: alert here
			ctxLogger.Error().Err(err).Msgf("%s: unable to find dictionary by dictionary_id", op)

			return c.Send(ui.DictionaryNotFoundMsg, ui.BuildMainMenuReplyKb())

		case errors.Is(err, domain.ErrSubscriptionNotFound):
			ctxLogger.Debug().Msgf("%s: not subscribed", op)

			return c.Send(ui.DictionarySubscriptionNotFoundMsg, ui.BuildMainMenuReplyKb())

		default:
			ctxLogger.Error().Err(err).Msgf("%s failed", op)

			return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
		}
	}

	ctxLogger.Debug().Msgf("%s confirmation requested", op)

	return c.Send(
		ui.DictionaryUnsubscribeConfirmMsg,
		&tele.SendOptions{ReplyMarkup: ui.BuildUnsubscribeConfirmInlineKb(dictionaryID)},
	)
}

func (h *BotHandlers) ConfirmUnsubscribe(c tele.Context) error {
	const op = "ConfirmUnsubscribe"

	ctx, cancel := context.WithTimeout(context.Background(), handlerCtxTimeout)
	defer cancel()

	userID := c.Sender().ID
	updateID := c.Update().ID
	username := c.Sender().Username

	// TODO: remove personal data from logs after alfa-test
	ctxLogger := h.logger.With().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Logger()

	if c.Callback() != nil {
		defer func() {
			_ = c.Respond()
			if err := c.Delete(); err != nil {
				ctxLogger.Warn().Err(err).Msgf("%s: failed to delete confirm message", op)
			}
		}()
	}

	ctxLogger.Debug().Msgf("handling %s", op)

	dictionaryID := extractCallbackDictionaryID(c)
	if dictionaryID == "" {
		// TODO: alert here
		ctxLogger.Error().Msgf("%s: dictionary_id is empty", op)

		return c.Send(ui.DictionaryNotFoundMsg, ui.BuildMainMenuReplyKb())
	}

	ctxLogger = ctxLogger.With().Str("dictionary_id", dictionaryID).Logger()

	err := h.subsUC.Unsubscribe(ctx, userID, dictionaryID)
	if err != nil {
		// TODO: move to mapper and use for ConfirmUnsubscribe & Unsubscribe
		// handlers
		switch {
		case errors.Is(err, domain.ErrDictionaryNotFound):
			// TODO: alert here
			ctxLogger.Error().Err(err).Msgf("%s: unable to find dictionary by dictionary_id", op)

			return c.Send(ui.DictionaryNotFoundMsg, ui.BuildMainMenuReplyKb())

		case errors.Is(err, domain.ErrSubscriptionNotFound):
			ctxLogger.Debug().Msgf("%s: not subscribed", op)

			return c.Send(ui.DictionarySubscriptionNotFoundMsg, ui.BuildMainMenuReplyKb())

		default:
			ctxLogger.Error().Err(err).Msgf("%s failed", op)

			return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
		}
	}

	ctxLogger.Debug().Msgf("%s handled", op)

	return c.Send(ui.DictionaryUnsubscribedMsg, ui.BuildMainMenuReplyKb())
}

func (h *BotHandlers) RejectUnsubscribe(c tele.Context) error {
	const op = "RejectUnsubscribe"

	// TODO (high): add basic handler && add there helper-funcs for doing staff
	// like init ctxLogger, getting request params, validation
	userID := c.Sender().ID
	updateID := c.Update().ID
	username := c.Sender().Username

	// TODO: remove personal data from logs after alfa-test
	ctxLogger := h.logger.With().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Logger()

	if c.Callback() != nil {
		defer func() {
			_ = c.Respond()
			if err := c.Delete(); err != nil {
				ctxLogger.Warn().Err(err).Msgf("%s: failed to delete confirm message", op)
			}
		}()
	}

	ctxLogger.Debug().Msgf("%s handled", op)

	return c.Send(ui.DictionaryUnsubscribeCanceledMsg, ui.BuildMainMenuReplyKb())
}

func (h *BotHandlers) LearnByDictNum(c tele.Context) error {
	const op = "LearnByDictNum"

	ctx, cancel := context.WithTimeout(context.Background(), handlerCtxTimeout)
	defer cancel()
	if c.Callback() != nil {
		defer func() {
			_ = c.Respond()
		}()
	}

	userID := c.Sender().ID
	updateID := c.Update().ID
	username := c.Sender().Username

	ctxLogger := h.logger.With().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Logger()

	ctxLogger.Debug().Msgf("handling %s", op)

	args := c.Args()
	if len(args) != 1 {
		ctxLogger.Debug().Int("args", len(args)).Msgf("%s: incorrect num of args", op)

		return c.Send(ui.LearnUsageMsg, ui.BuildMainMenuReplyKb())
	}

	trimmed := strings.Trim(strings.TrimSpace(args[0]), "<>")
	number, convErr := strconv.Atoi(trimmed)
	if convErr != nil {
		ctxLogger.Debug().
			Err(convErr).
			Str("args[0]", args[0]).
			Msgf("%s: error converting arg to int", op)

		return c.Send(ui.LearnUsageMsg, ui.BuildMainMenuReplyKb())
	}

	word, dictionaryID, err := h.learnUC.LearnByDictionaryNumber(ctx, userID, number)
	if err != nil {
		mapped := mapper.MapLearningErrorToUI(err)
		if mapped.State() != mapper.LearningUIUnknown {
			ctxLogger.Debug().
				Err(err).
				Str("dictionary_id", dictionaryID).
				Msgf("%s handled with mapped learning error", op)

			return mapper.SendLearningMappedError(c, mapped, dictionaryID)
		}

		// TODO: alert here
		ctxLogger.Error().Err(err).Msgf("%s failed", op)

		return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
	}

	return c.Send(
		ui.FormatLearningWordCard(*word),
		&tele.SendOptions{
			ParseMode:   tele.ModeHTML,
			ReplyMarkup: ui.BuildLearningReplyKb(),
		},
	)
}

func (h *BotHandlers) LearnByDictID(c tele.Context) error {
	const op = "LearnByDictID"

	ctx, cancel := context.WithTimeout(context.Background(), handlerCtxTimeout)
	defer cancel()
	if c.Callback() != nil {
		defer func() {
			_ = c.Respond()
		}()
	}

	userID := c.Sender().ID
	updateID := c.Update().ID
	username := c.Sender().Username

	ctxLogger := h.logger.With().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Logger()

	ctxLogger.Debug().Msgf("handling %s", op)

	dictionaryID := extractCallbackDictionaryID(c)
	if dictionaryID == "" {
		// TODO: alert here
		ctxLogger.Error().Msgf("%s: dictionary_id is empty", op)

		return c.Send(ui.DictionaryNotFoundMsg, ui.BuildMainMenuReplyKb())
	}

	word, err := h.learnUC.LearnByDictionaryID(ctx, userID, dictionaryID)
	if err != nil {
		mapped := mapper.MapLearningErrorToUI(err)
		if mapped.State() != mapper.LearningUIUnknown {
			ctxLogger.Debug().
				Err(err).
				Str("dictionary_id", dictionaryID).
				Msgf("%s handled with mapped learning error", op)

			return mapper.SendLearningMappedError(c, mapped, dictionaryID)
		}

		ctxLogger.Error().Err(err).Msgf("%s failed", op)

		return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
	}

	return c.Send(
		ui.FormatLearningWordCard(*word),
		&tele.SendOptions{
			ParseMode:   tele.ModeHTML,
			ReplyMarkup: ui.BuildLearningReplyKb(),
		},
	)
}

func (h *BotHandlers) LearningAction(c tele.Context) error {
	const op = "LearningAction"

	ctx, cancel := context.WithTimeout(context.Background(), handlerCtxTimeout)
	defer cancel()
	if c.Callback() != nil {
		defer func() {
			_ = c.Respond()
		}()
	}

	userID := c.Sender().ID
	updateID := c.Update().ID
	username := c.Sender().Username

	ctxLogger := h.logger.With().
		Int("update_id", updateID).
		Int64("user_id", userID).
		Str("username", username).
		Logger()

	ctxLogger.Debug().Msgf("handling %s", op)

	switch c.Text() {
	case ui.LearnAddText:
		return h.handleLearningDecision(ctx, userID, h.learnUC.AddCurrentWord, ctxLogger, op, c)
	case ui.LearnBlockText:
		return h.handleLearningDecision(ctx, userID, h.learnUC.BlockCurrentWord, ctxLogger, op, c)
	case ui.LearnReviewText:
		return c.Send(ui.LearnReviewWIPMsg, ui.BuildLearningReplyKb())
	case ui.LearnBackText:
		if err := h.learnUC.Back(ctx, userID); err != nil {
			ctxLogger.Error().Err(err).Msgf("%s failed", op)

			return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
		}

		return c.Send(ui.ToMainMenuMsg, ui.BuildMainMenuReplyKb())
	default:
		return nil
	}
}

func (h *BotHandlers) handleLearningDecision(
	ctx context.Context,
	userID int64,
	decisionFn func(context.Context, int64) (*domain.LearningWord, error),
	ctxLogger zerolog.Logger,
	op string,
	c tele.Context,
) error {
	word, err := decisionFn(ctx, userID)
	if err != nil {
		dictionaryID, _ := h.learnUC.ActiveDictionaryID(ctx, userID)

		mapped := mapper.MapLearningErrorToUI(err)
		if mapped.State() != mapper.LearningUIUnknown {
			ctxLogger.Debug().
				Err(err).
				Str("dictionary_id", dictionaryID).
				Msgf("%s handled with mapped learning error", op)

			return mapper.SendLearningMappedError(c, mapped, dictionaryID)
		}

		ctxLogger.Error().Err(err).Msgf("%s failed", op)

		return c.Send(ui.InternalErrorMsg, ui.BuildMainMenuReplyKb())
	}

	return c.Send(
		ui.FormatLearningWordCard(*word),
		&tele.SendOptions{
			ParseMode:   tele.ModeHTML,
			ReplyMarkup: ui.BuildLearningReplyKb(),
		},
	)
}

func (h *BotHandlers) Review(c tele.Context) error {
	if c.Callback() != nil {
		defer func() {
			_ = c.Respond()
		}()
	}

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

func extractCallbackDictionaryID(c tele.Context) string {
	if c.Callback() != nil {
		return strings.TrimSpace(c.Data())
	}

	return ""
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
