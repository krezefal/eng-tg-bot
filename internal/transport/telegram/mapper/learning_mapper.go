package mapper

import (
	"errors"

	tele "gopkg.in/telebot.v4"

	"github.com/krezefal/eng-tg-bot/internal/domain"
	"github.com/krezefal/eng-tg-bot/internal/transport/telegram/ui"
)

type LearningUIState int

const (
	LearningUIUnknown LearningUIState = iota
	LearningUIMainMenu
	LearningUICompleted
)

type LearningUIResult struct {
	state LearningUIState
	msg   string
}

func (lr LearningUIResult) State() LearningUIState {
	return lr.state
}

func (lr LearningUIResult) Message() string {
	return lr.msg
}

func MapLearningErrorToUI(err error) LearningUIResult {
	switch {
	case errors.Is(err, domain.ErrInvalidDictionaryNumber):
		return LearningUIResult{state: LearningUIMainMenu, msg: ui.LearnInvalidDictionaryNumberMsg}
	case errors.Is(err, domain.ErrDictionaryNotFound):
		return LearningUIResult{state: LearningUIMainMenu, msg: ui.DictionaryNotFoundMsg}
	case errors.Is(err, domain.ErrSubscriptionNotFound):
		return LearningUIResult{state: LearningUIMainMenu, msg: ui.LearnNotSubscribedMsg}
	case errors.Is(err, domain.ErrLearningNotStarted):
		return LearningUIResult{state: LearningUIMainMenu, msg: ui.LearnNotStartedMsg}
	case errors.Is(err, domain.ErrNoWordsForLearning):
		return LearningUIResult{state: LearningUICompleted, msg: ui.LearnCompletedMsg}
	default:
		return LearningUIResult{state: LearningUIUnknown}
	}
}

func SendLearningMappedError(c tele.Context, mapped LearningUIResult, dictionaryID string) error {
	switch mapped.state {
	case LearningUIMainMenu:
		return c.Send(mapped.msg, ui.BuildMainMenuReplyKb())
	case LearningUICompleted:
		return c.Send(
			mapped.msg,
			&tele.SendOptions{ReplyMarkup: ui.BuildLearningCompletedInlineKb(dictionaryID)},
		)
	default:
		return nil
	}
}
