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
	b.WriteString(fmt.Sprintf("📘 <u>%s</u>\n", html.EscapeString(title)))

	if strings.TrimSpace(dict.Description) != "" {
		b.WriteString(fmt.Sprintf("Описание: %s\n", html.EscapeString(dict.Description)))
	}

	if strings.TrimSpace(dict.Author) != "" {
		b.WriteString(fmt.Sprintf("Автор: %s\n", html.EscapeString(dict.Author)))
	}

	b.WriteString(fmt.Sprintf("Тип: %s", html.EscapeString(dict.Mode.HumanReadable())))

	return b.String()
}

func FormatSubscribedDictionaryCard(number int, dict domain.Dictionary) string {
	var b strings.Builder
	title := strings.TrimSpace(dict.Title)
	if title == "" {
		title = "Без названия"
	}
	b.WriteString(fmt.Sprintf("%d. 📘 <u>%s</u>\n", number, html.EscapeString(title)))

	if strings.TrimSpace(dict.Author) != "" {
		b.WriteString(fmt.Sprintf("Автор: %s\n", html.EscapeString(dict.Author)))
	}

	b.WriteString(fmt.Sprintf("Тип: %s", html.EscapeString(dict.Mode.HumanReadable())))

	return b.String()
}

func FormatDictionaryDetails(dict domain.Dictionary, words []domain.DictionaryWordPreview) string {
	var b strings.Builder
	title := strings.TrimSpace(dict.Title)
	if title == "" {
		title = "Без названия"
	}
	b.WriteString(fmt.Sprintf("📘 <u>%s</u>\n\n", html.EscapeString(title)))

	if strings.TrimSpace(dict.Description) != "" {
		b.WriteString(fmt.Sprintf("Описание: %s\n", html.EscapeString(dict.Description)))
	}

	if strings.TrimSpace(dict.Author) != "" {
		b.WriteString(fmt.Sprintf("Автор: %s\n", html.EscapeString(dict.Author)))
	}

	dictModeHint := ""
	switch dict.Mode {
	case domain.RandomPoolMode:
		dictModeHint = "ты сам выбираешь, когда приступать к изучению новых слов"
	case domain.OnScheduleMode:
		dictModeHint = "новые слова приходят тебе по расписанию, заданному автором"
	}

	b.WriteString(fmt.Sprintf("Тип: %s — %s\n\n",
		html.EscapeString(dict.Mode.HumanReadable()), html.EscapeString(dictModeHint)))

	if len(words) == 0 {
		b.WriteString("Примеры слов: пока нет слов в словаре")

		return b.String()
	}

	b.WriteString("Примеры слов:\n")
	for _, w := range words {
		b.WriteString(
			fmt.Sprintf("• %s — <tg-spoiler>%s</tg-spoiler>\n",
				html.EscapeString(w.Spelling), html.EscapeString(w.RUTranslation)),
		)
	}

	return strings.TrimSpace(b.String())
}

func FormatLearningWordCard(word domain.LearningWord) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("🇬🇧 <b>%s</b> — %s\n\n",
		html.EscapeString(word.Spelling), html.EscapeString(word.Transcription)))

	b.WriteString(fmt.Sprintf("🇷🇺 <tg-spoiler>%s</tg-spoiler>\n\n", html.EscapeString(word.RUTranslation)))

	b.WriteString("Примеры использования в речи:\n" +
		"*тут может быть куча предложений, чтобы лучше понять контекст употребления, но пока их нет*")

	return b.String()
}
