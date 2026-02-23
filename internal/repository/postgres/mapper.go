package postgres

import (
	"fmt"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

type rowScanner interface {
	Scan(dest ...any) error
}

func toDomainDictionary(scanner rowScanner) (*domain.Dictionary, error) {
	var d domain.Dictionary
	var rawMode string
	err := scanner.Scan(&d.ID, &d.Title, &d.Description, &rawMode, &d.Author, &d.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to convert into dictionary: %w", err)
	}

	mode, ok := domain.ParseDictionaryMode(rawMode)
	if !ok {
		return nil, fmt.Errorf("unsupported dictionary mode: %q", rawMode)
	}
	d.Mode = mode

	return &d, nil
}

func toDomainDictionaryWordPreview(scanner rowScanner) (*domain.DictionaryWordPreview, error) {
	var w domain.DictionaryWordPreview
	err := scanner.Scan(&w.Spelling, &w.RUTranslation)
	if err != nil {
		return nil, fmt.Errorf("failed to convert into dictionary word preview: %w", err)
	}

	return &w, nil
}

func toDomainLearningWord(scanner rowScanner) (*domain.LearningWord, error) {
	var w domain.LearningWord
	err := scanner.Scan(
		&w.ID,
		&w.DictionaryID,
		&w.Spelling,
		&w.Transcription,
		&w.Audio,
		&w.RUTranslation,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to convert into learning word: %w", err)
	}

	return &w, nil
}
