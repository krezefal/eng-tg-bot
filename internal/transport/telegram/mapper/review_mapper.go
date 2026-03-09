package mapper

import (
	"errors"

	"github.com/krezefal/eng-tg-bot/internal/domain"
	"github.com/krezefal/eng-tg-bot/internal/transport/telegram/ui"
	tele "gopkg.in/telebot.v4"
)

type ReviewUIState int

const (
	ReviewUIUnknown ReviewUIState = iota
	ReviewUIMainMenu
	ReviewUINoDue
	ReviewUIDone
)

var ReviewGradeByText = map[string]int{
	ui.ReviewRate1Text: 0,
	ui.ReviewRate2Text: 1,
	ui.ReviewRate3Text: 2,
	ui.ReviewRate4Text: 3,
}

type ReviewUIResult struct {
	state ReviewUIState
	msg   string
}

func (rr ReviewUIResult) State() ReviewUIState {
	return rr.state
}

// TODO: pass logger here to prevent leak of business logic decisions to
// transport layer -> decisions about loggin important errors.
func MapReviewErrorToUI(err error) *ReviewUIResult {
	switch {
	case errors.Is(err, domain.ErrInvalidDictionaryNumber):
		return &ReviewUIResult{state: ReviewUIMainMenu, msg: ui.InvalidDictionaryNumberMsg}
	case errors.Is(err, domain.ErrDictionaryNotFound):
		return &ReviewUIResult{state: ReviewUIMainMenu, msg: ui.DictionaryNotFoundMsg}
	case errors.Is(err, domain.ErrSubscriptionNotFound):
		return &ReviewUIResult{state: ReviewUIMainMenu, msg: ui.NotSubscribedMsg}
	case errors.Is(err, domain.ErrReviewNotStarted):
		// TODO: logger + alert here
		return &ReviewUIResult{state: ReviewUIMainMenu, msg: ui.ActiveDictMissingMsg}
	case errors.Is(err, domain.ErrEmptyReviewWordsList):
		return &ReviewUIResult{state: ReviewUIMainMenu, msg: ui.ReviewEmptyWordsListMsg}
	case errors.Is(err, domain.ErrNoWordsDueForReview):
		return &ReviewUIResult{state: ReviewUINoDue, msg: ui.ReviewNoDueMsg}
	case errors.Is(err, domain.ErrReviewRoundFinished):
		return &ReviewUIResult{state: ReviewUIDone, msg: ui.ReviewCompletedMsg}
	default:
		return &ReviewUIResult{state: ReviewUIUnknown}
	}
}

func SendReviewMappedError(c tele.Context, mapped *ReviewUIResult, dictionaryID string) error {
	switch mapped.state {
	case ReviewUIMainMenu:
		return c.Send(mapped.msg, ui.BuildMainMenuReplyKb())
	case ReviewUINoDue:
		if dictionaryID == "" {
			return c.Send(mapped.msg, ui.BuildMainMenuReplyKb())
		}

		return c.Send(
			mapped.msg,
			&tele.SendOptions{ReplyMarkup: ui.BuildReviewForceInlineKb(dictionaryID)},
		)
	case ReviewUIDone:
		return c.Send(mapped.msg, ui.BuildReviewFinishReplyKb())
	default:
		return nil
	}
}
