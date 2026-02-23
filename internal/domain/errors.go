package domain

import "errors"

var (
	ErrDictionaryNotFound   = errors.New("dictionary not found")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrAlreadySubscribed    = errors.New("already subscribed")
)
