package domain

type UserWordStatus string

const (
	UserWordStatusLearning UserWordStatus = "learning"
	UserWordStatusBlocked  UserWordStatus = "blocked"
)
