package domain

import "time"

type DictionaryMode int

const (
	UnsupportedMode DictionaryMode = iota
	RandomPoolMode
	OnScheduleMode
)

func (m DictionaryMode) HumanReadable() string {
	switch m {
	case RandomPoolMode:
		return "Обычный словарь"
	case OnScheduleMode:
		return "Словарь-по-расписанию"
	default:
		return "unknown"
	}
}

func ParseDictionaryMode(raw string) (DictionaryMode, bool) {
	switch raw {
	case "random_pool":
		return RandomPoolMode, true
	case "on_schedule":
		return OnScheduleMode, true
	default:
		return UnsupportedMode, false
	}
}

type Dictionary struct {
	ID          string
	Title       string
	Description string
	Mode        DictionaryMode
	Author      string
	CreatedAt   time.Time
}

type DictionaryDetails struct {
	Dictionary *Dictionary
	Words      []DictionaryWordPreview
}
