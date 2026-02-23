package domain

import "errors"

var (
	ErrDictionaryNotFound      = errors.New("dictionary not found")
	ErrSubscriptionNotFound    = errors.New("subscription not found")
	ErrAlreadySubscribed       = errors.New("already subscribed")
	ErrInvalidDictionaryNumber = errors.New("invalid dictionary number")
	ErrNoWordsForLearning      = errors.New("no words for learning")
	ErrLearningNotStarted      = errors.New("learning not started")
)
