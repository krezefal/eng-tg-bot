package ui

import (
	"fmt"
	"html"
	"strings"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

func FormatDictionaryCard(dict domain.Dictionary) string {
	var b strings.Builder
	title := strings.TrimSpace(dict.Title)
	if title == "" {
		title = "Без названия"
	}
	fmt.Fprintf(&b, "📘 <u>%s</u>\n", html.EscapeString(title))

	if strings.TrimSpace(dict.Description) != "" {
		fmt.Fprintf(&b, "Описание: %s\n", html.EscapeString(dict.Description))
	}

	if strings.TrimSpace(dict.Author) != "" {
		fmt.Fprintf(&b, "Автор: %s\n", html.EscapeString(dict.Author))
	}

	fmt.Fprintf(&b, "Тип: %s", html.EscapeString(dict.Mode.HumanReadable()))

	return b.String()
}

func FormatSubscribedDictionaryCard(number int, dict domain.Dictionary) string {
	var b strings.Builder
	title := strings.TrimSpace(dict.Title)
	if title == "" {
		title = "Без названия"
	}
	fmt.Fprintf(&b, "%d. 📘 <u>%s</u>\n", number, html.EscapeString(title))

	if strings.TrimSpace(dict.Author) != "" {
		fmt.Fprintf(&b, "Автор: %s\n", html.EscapeString(dict.Author))
	}

	fmt.Fprintf(&b, "Тип: %s", html.EscapeString(dict.Mode.HumanReadable()))

	return b.String()
}

func FormatDictionaryDetails(dict domain.Dictionary, words []domain.DictionaryWordPreview) string {
	var b strings.Builder
	title := strings.TrimSpace(dict.Title)
	if title == "" {
		title = "Без названия"
	}
	fmt.Fprintf(&b, "📘 <u>%s</u>\n\n", html.EscapeString(title))

	if strings.TrimSpace(dict.Description) != "" {
		fmt.Fprintf(&b, "Описание: %s\n", html.EscapeString(dict.Description))
	}

	if strings.TrimSpace(dict.Author) != "" {
		fmt.Fprintf(&b, "Автор: %s\n", html.EscapeString(dict.Author))
	}

	dictModeHint := ""
	switch dict.Mode {
	case domain.RandomPoolMode:
		dictModeHint = "ты сам выбираешь, когда приступать к изучению новых слов"
	case domain.OnScheduleMode:
		dictModeHint = "новые слова приходят тебе по расписанию, заданному автором"
	}

	fmt.Fprintf(&b, "Тип: %s — %s\n\n",
		html.EscapeString(dict.Mode.HumanReadable()), html.EscapeString(dictModeHint))

	if len(words) == 0 {
		b.WriteString("В этом словаре пока нет слов 💤")

		return b.String()
	}

	b.WriteString("Несколько слов отсюда:\n")
	for _, w := range words {
		fmt.Fprintf(&b, "• %s — <tg-spoiler>%s</tg-spoiler>\n",
			html.EscapeString(w.Spelling), html.EscapeString(w.RUTranslation))
	}

	return strings.TrimSpace(b.String())
}

func FormatLearningWordCard(word domain.LearningWord) string {
	var b strings.Builder
	fmt.Fprintf(&b, "🇬🇧 <b>%s</b> — %s\n\n",
		html.EscapeString(word.Spelling), html.EscapeString(word.Transcription))

	fmt.Fprintf(&b, "🇷🇺 <tg-spoiler>%s</tg-spoiler>\n\n", html.EscapeString(word.RUTranslation))

	b.WriteString("Примеры использования в речи:\n" +
		"*тут может быть куча предложений, чтобы лучше понять контекст употребления, но пока их нет*")

	return b.String()
}

func FormatReviewWordCard(word *domain.ReviewWord) string {
	var b strings.Builder
	fmt.Fprintf(&b, "🇬🇧 <b>%s</b> — %s\n\n",
		html.EscapeString(word.Spelling), html.EscapeString(word.Transcription))
	fmt.Fprintf(&b, "🇷🇺 <tg-spoiler>%s</tg-spoiler>", html.EscapeString(word.RUTranslation))

	return b.String()
}
