package ui

import tele "gopkg.in/telebot.v4"

const (
	MainMenuDictText   = "📚 Словари (публичные)"
	MainMenuMyDictText = "📖 Мои словари"
	MainMenuHelpText   = "❔ Помощь"
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

	btnSubscribe := markup.Data("Подписаться", "dict_subscribe", dictionaryID)
	btnDetails := markup.Data("Подробнее", "dict_details", dictionaryID)

	markup.Inline(
		markup.Row(btnSubscribe, btnDetails),
	)

	return markup
}

func BuildUserDictionaryInlineKb(dictionaryID string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnLearn := markup.Data("Учить", "dict_learn", dictionaryID)
	btnReview := markup.Data("Повторить", "dict_review", dictionaryID)
	btnUnsubscribe := markup.Data("Отписаться", "dict_unsubscribe", dictionaryID)

	markup.Inline(
		markup.Row(btnLearn),
		markup.Row(btnReview, btnUnsubscribe),
	)

	return markup
}

func BuildDictionaryDetailsInlineKb(dictionaryID string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnSubscribe := markup.Data("Подписаться", "dict_subscribe", dictionaryID)
	btnDetails := markup.Data("К словарям", "to_dicts")

	markup.Inline(
		markup.Row(btnSubscribe, btnDetails),
	)

	return markup
}

func BuildUnsubscribeConfirmInlineKb(dictionaryID string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnConfirm := markup.Data("Да", "dict_confirm_unsubscribe", dictionaryID)
	btnReject := markup.Data("Нет", "dict_reject_unsubscribe", dictionaryID)

	markup.Inline(
		markup.Row(btnConfirm, btnReject),
	)

	return markup
}
