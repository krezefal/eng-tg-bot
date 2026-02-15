package domain

import tele "gopkg.in/telebot.v4"

func BuildRateKeyboard(wordID string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btn1 := markup.Data("Не помню", "rate", wordID+":1")
	btn2 := markup.Data("Слабо помню", "rate", wordID+":2")
	btn3 := markup.Data("Возможно вспомнил бы", "rate", wordID+":3")
	btn4 := markup.Data("Хорошо помню", "rate", wordID+":4")
	btn5 := markup.Data("Запомнил!", "rate", wordID+":5")

	markup.Inline(
		markup.Row(btn5),
		markup.Row(btn3, btn4),
		markup.Row(btn1, btn2),
	)

	return markup
}

func BuildLearningKeyboard(wordID string) *tele.ReplyMarkup {
	markup := &tele.ReplyMarkup{}

	btnAdd := markup.Data("➕ Добавить", "learn", wordID+":add")
	btnSkip := markup.Data("⏭ Пропустить", "learn", wordID+":skip")
	btnBlock := markup.Data("🚫 Не предлагать", "learn", wordID+":block")

	markup.Inline(
		markup.Row(btnAdd),
		markup.Row(btnSkip),
		markup.Row(btnBlock),
	)

	return markup
}
