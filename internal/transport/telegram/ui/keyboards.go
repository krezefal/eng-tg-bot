package ui

import tele "gopkg.in/telebot.v4"

// Only reply btns contain emoji.
const (
	MainMenuDictText   = "📚 Словари (публичные)"
	MainMenuMyDictText = "📖 Мои словари"
	MainMenuHelpText   = "❔ Помощь"

	AddDictText     = "Добавить словарь"
	DictDetailsText = "Подробнее"
	ToDictsText     = "К словарям"
	RemoveDictText  = "Отписаться"

	ConfirnUnsubText = "Да"
	RejectUnsubText  = "Нет"

	StartLearnText  = "Учить"
	StartReviewText = "Повторить"

	LearnAddText       = "✍️ Добавить в словарь"
	LearnBlockText     = "🙅‍♂️ Не добавлять — знаю это слово"
	LearnReviewText    = "🧠 Перейти к повторению слов"
	LearnReviewNowText = "Перейти к повторению сейчас"

	ReviewStartText   = "🚀 Старт"
	ReviewRestartText = "🔁 Повторить еще раз"
	ReviewStopText    = "️🏁 Закончить подход"
	ReviewRate1Text   = "Не помню"
	ReviewRate2Text   = "Трудно"
	ReviewRate3Text   = "Легко"
	ReviewRate4Text   = "Помню!"
	ReviewForceStart  = "Все равно хочу попрактиковаться"

	ToMainMenuText = "🏠 В главное меню"
)

func BuildMainMenuReplyKb() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnDict := markup.Text(MainMenuDictText)
	btnList := markup.Text(MainMenuMyDictText)
	btnHelp := markup.Text(MainMenuHelpText)

	markup.Reply(
		markup.Row(btnDict),
		markup.Row(btnList),
		markup.Row(btnHelp),
	)

	return markup
}

func BuildPublicDictionaryInlineKb(dictionaryID string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnSubscribe := markup.Data(AddDictText, "dict_subscribe", dictionaryID)
	btnDetails := markup.Data(DictDetailsText, "dict_details", dictionaryID)

	markup.Inline(
		markup.Row(btnSubscribe, btnDetails),
	)

	return markup
}

func BuildUserDictionaryInlineKb(dictionaryID string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnLearn := markup.Data(StartLearnText, "dict_learn", dictionaryID)
	btnReview := markup.Data(StartReviewText, "dict_review", dictionaryID)
	btnUnsubscribe := markup.Data(RemoveDictText, "dict_unsubscribe", dictionaryID)

	markup.Inline(
		markup.Row(btnLearn),
		markup.Row(btnReview, btnUnsubscribe),
	)

	return markup
}

func BuildDictionaryDetailsInlineKb(dictionaryID string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnSubscribe := markup.Data(AddDictText, "dict_subscribe", dictionaryID)
	btnDetails := markup.Data(ToDictsText, "to_dicts")

	markup.Inline(
		markup.Row(btnSubscribe, btnDetails),
	)

	return markup
}

func BuildUnsubscribeConfirmInlineKb(dictionaryID string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnConfirm := markup.Data(ConfirnUnsubText, "dict_confirm_unsubscribe", dictionaryID)
	btnReject := markup.Data(RejectUnsubText, "dict_reject_unsubscribe", dictionaryID)

	markup.Inline(
		markup.Row(btnConfirm, btnReject),
	)

	return markup
}

func BuildLearningReplyKb() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnAdd := markup.Text(LearnAddText)
	btnBlock := markup.Text(LearnBlockText)
	btnReview := markup.Text(LearnReviewText)
	btnBack := markup.Text(ToMainMenuText)

	markup.Reply(
		markup.Row(btnAdd, btnBlock),
		markup.Row(btnReview),
		markup.Row(btnBack),
	)

	return markup
}

func BuildLearningCompletedInlineKb(dictionaryID string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnReview := markup.Data(LearnReviewNowText, "dict_review", dictionaryID)
	markup.Inline(markup.Row(btnReview))

	return markup
}

func BuildLearningCompletedReplyKb() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnReview := markup.Text(LearnReviewText)
	btnBack := markup.Text(ToMainMenuText)

	markup.Reply(
		markup.Row(btnReview),
		markup.Row(btnBack),
	)

	return markup
}

func BuildReviewIntroReplyKb() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnStart := markup.Text(ReviewStartText)
	btnMain := markup.Text(ToMainMenuText)

	markup.Reply(
		markup.Row(btnStart),
		markup.Row(btnMain),
	)

	return markup
}

func BuildReviewRateReplyKb() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{ResizeKeyboard: true}

	btn1 := markup.Text(ReviewRate1Text)
	btn2 := markup.Text(ReviewRate2Text)
	btn3 := markup.Text(ReviewRate3Text)
	btn4 := markup.Text(ReviewRate4Text)
	btnStop := markup.Text(ReviewStopText)

	markup.Reply(
		markup.Row(btn1, btn2, btn3, btn4),
		markup.Row(btnStop),
	)

	return markup
}

func BuildReviewForceInlineKb(dictionaryID string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnForce := markup.Data(ReviewForceStart, "review_force", dictionaryID)
	markup.Inline(markup.Row(btnForce))

	return markup
}

func BuildReviewFinishReplyKb() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnRestart := markup.Text(ReviewRestartText)
	btnMain := markup.Text(ToMainMenuText)

	markup.Reply(
		markup.Row(btnRestart),
		markup.Row(btnMain),
	)

	return markup
}
