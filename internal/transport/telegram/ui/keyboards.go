package ui

import tele "gopkg.in/telebot.v4"

const (
	MainMenuDictText   = "📚 Словари (публичные)"
	MainMenuMyDictText = "⭐ Мои словари"
	MainMenuHelpText   = "❓ Помощь"
)

func BuildMainMenuKeyboard() *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnDict := markup.Text(MainMenuDictText)
	btnMyDict := markup.Text(MainMenuMyDictText)
	btnHelp := markup.Text(MainMenuHelpText)

	markup.Reply(
		markup.Row(btnDict),
		markup.Row(btnMyDict),
		markup.Row(btnHelp),
	)

	return markup
}

func BuildRateKeyboard(wordID string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btn1 := markup.Data("Не помню", "rate", wordID+":1")
	btn2 := markup.Data("Слабо помню", "rate", wordID+":2")
	btn3 := markup.Data("Хорошо помню", "rate", wordID+":3")
	btn4 := markup.Data("Запомнил!", "rate", wordID+":4")

	btnStop := markup.Data("Закончили подход", "rate")

	markup.Inline(
		markup.Row(btn4, btn3),
		markup.Row(btn2, btn1),
		markup.Row(btnStop),
	)

	return markup
}

func BuildLearningKeyboard(wordID string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnAdd := markup.Data("➕ Добавить", "learn", wordID+":add")
	btnBlock := markup.Data("⏭ Знаю! Не предлагать", "learn", wordID+":block")

	btnStop := markup.Data("Закончили подход", "rate")

	markup.Inline(
		markup.Row(btnAdd),
		markup.Row(btnBlock),
		markup.Row(btnStop),
	)

	return markup
}
