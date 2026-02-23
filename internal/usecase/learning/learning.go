package learning

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

type Usecase struct {
	userRepo      UserRepo
	dictRepo      DictionaryRepo
	subsRepo      SubscriptionsRepo
	wordStateRepo WordStateRepo
	logger        *zerolog.Logger

	pendingMu   sync.RWMutex
	pendingWord map[int64]pendingWord
}

type pendingWord struct {
	dictionaryID string
	wordID       string
}

func NewUsecase(
	userRepo UserRepo,
	dictRepo DictionaryRepo,
	subsRepo SubscriptionsRepo,
	wordStateRepo WordStateRepo,
	parentLogger *zerolog.Logger,
) *Usecase {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}

	logger := parentLogger.With().Str("component", "learning_usecase").Logger()

	return &Usecase{
		userRepo:      userRepo,
		dictRepo:      dictRepo,
		subsRepo:      subsRepo,
		wordStateRepo: wordStateRepo,
		logger:        &logger,
		pendingWord:   make(map[int64]pendingWord),
	}
}

func (u *Usecase) LearnByDictionaryNumber(
	ctx context.Context,
	userID int64,
	number int,
) (*domain.LearningWord, string, error) {
	const op = "LearnByDictionaryNumber"

	if number <= 0 {
		return nil, "", domain.ErrInvalidDictionaryNumber
	}

	dictionaries, err := u.subsRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, "", fmt.Errorf("%s failed: %w", op, err)
	}

	if number > len(dictionaries) {
		return nil, "", domain.ErrInvalidDictionaryNumber
	}

	dictionaryID := dictionaries[number-1].ID
	word, err := u.startLearning(ctx, userID, dictionaryID)
	if err != nil {
		return nil, dictionaryID, fmt.Errorf("%s failed: %w", op, err)
	}

	u.logger.Debug().
		Int64("user_id", userID).
		Str("dictionary_id", dictionaryID).
		Str("word_id", word.ID).
		Msgf("%s succeeded", op)

	return word, dictionaryID, nil
}

func (u *Usecase) LearnByDictionaryID(
	ctx context.Context,
	userID int64,
	dictionaryID string,
) (*domain.LearningWord, error) {
	const op = "LearnByDictionaryID"

	exists, err := u.dictRepo.ExistsByID(ctx, dictionaryID)
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}
	if !exists {
		return nil, domain.ErrDictionaryNotFound
	}

	subscribed, err := u.subsRepo.IsSubscribedByUser(ctx, userID, dictionaryID)
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}
	if !subscribed {
		return nil, domain.ErrSubscriptionNotFound
	}

	word, err := u.startLearning(ctx, userID, dictionaryID)
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}

	u.logger.Debug().
		Int64("user_id", userID).
		Str("dictionary_id", dictionaryID).
		Str("word_id", word.ID).
		Msgf("%s succeeded", op)

	return word, nil
}

func (u *Usecase) startLearning(
	ctx context.Context,
	userID int64,
	dictionaryID string,
) (*domain.LearningWord, error) {
	const op = "startLearning"

	word, err := u.dictRepo.PickRandomUntrackedWord(ctx, userID, dictionaryID)
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}
	if word == nil {
		u.clearPending(userID)
		return nil, domain.ErrNoWordsForLearning
	}

	u.setPending(userID, pendingWord{
		dictionaryID: dictionaryID,
		wordID:       word.ID,
	})

	if err = u.subsRepo.MarkLearningStarted(ctx, userID, dictionaryID); err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}
	if err = u.userRepo.SetActiveDictionaryID(ctx, userID, dictionaryID); err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}

	return word, nil
}

func (u *Usecase) AddCurrentWord(ctx context.Context, userID int64) (*domain.LearningWord, error) {
	const op = "AddCurrentWord"

	return u.applyDecision(ctx, userID, domain.UserWordStatusLearning, op)
}

func (u *Usecase) BlockCurrentWord(ctx context.Context, userID int64) (*domain.LearningWord, error) {
	const op = "BlockCurrentWord"

	return u.applyDecision(ctx, userID, domain.UserWordStatusBlocked, op)
}

func (u *Usecase) ActiveDictionaryID(ctx context.Context, userID int64) (string, error) {
	const op = "ActiveDictionaryID"

	dictionaryID, err := u.userRepo.GetActiveDictionaryID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s failed: %w", op, err)
	}

	return dictionaryID, nil
}

func (u *Usecase) Back(ctx context.Context, userID int64) error {
	const op = "Back"

	u.clearPending(userID)
	if err := u.userRepo.ClearActiveDictionaryID(ctx, userID); err != nil {
		return fmt.Errorf("%s failed: %w", op, err)
	}

	u.logger.Debug().
		Int64("user_id", userID).
		Msgf("%s succeeded", op)

	return nil
}

func (u *Usecase) applyDecision(
	ctx context.Context,
	userID int64,
	status domain.UserWordStatus,
	op string,
) (*domain.LearningWord, error) {
	current, ok := u.getPending(userID)
	if !ok {
		return nil, domain.ErrLearningNotStarted
	}

	// TODO (high priority): set default ef, last_reviewed_at, etc - check how
	// it was implemented in 1st version of bot
	if err := u.wordStateRepo.UpsertStatus(ctx, userID, current.wordID, status); err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}

	nextWord, err := u.dictRepo.PickRandomUntrackedWord(ctx, userID, current.dictionaryID)
	if err != nil {
		return nil, fmt.Errorf("%s failed: %w", op, err)
	}
	if nextWord == nil {
		u.clearPending(userID)
		return nil, domain.ErrNoWordsForLearning
	}

	u.setPending(userID, pendingWord{
		dictionaryID: current.dictionaryID,
		wordID:       nextWord.ID,
	})

	u.logger.Debug().
		Int64("user_id", userID).
		Str("dictionary_id", current.dictionaryID).
		Str("word_id", nextWord.ID).
		Str("decision", string(status)).
		Msgf("%s succeeded", op)

	return nextWord, nil
}

func (u *Usecase) setPending(userID int64, p pendingWord) {
	u.pendingMu.Lock()
	defer u.pendingMu.Unlock()
	u.pendingWord[userID] = p
}

func (u *Usecase) getPending(userID int64) (pendingWord, bool) {
	u.pendingMu.RLock()
	defer u.pendingMu.RUnlock()
	p, ok := u.pendingWord[userID]

	return p, ok
}

func (u *Usecase) clearPending(userID int64) {
	u.pendingMu.Lock()
	defer u.pendingMu.Unlock()
	delete(u.pendingWord, userID)
}
