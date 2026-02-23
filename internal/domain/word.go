package domain

type DictionaryWordPreview struct {
	Spelling      string
	RUTranslation string
}

type LearningWord struct {
	ID            string
	DictionaryID  string
	Spelling      string
	Transcription string
	Audio         string
	RUTranslation string
}
