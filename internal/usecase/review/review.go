package review

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/krezefal/eng-tg-bot/internal/domain"
)

type Usecase struct {
	dictionaryRepo DictionaryRepo
	subsRepo       SubscriptionsRepo
	wordStateRepo  WordsStateRepo
	logger         *zerolog.Logger

	sessionMu sync.RWMutex
	sessions  map[int64]*reviewSession
}

type reviewSession struct {
	dictionaryID string
	queue        []*domain.ReviewWord
	current      *domain.ReviewWord
}

func NewUsecase(
	dictionaryRepo DictionaryRepo,
	subsRepo SubscriptionsRepo,
	wordStateRepo WordsStateRepo,
	parentLogger *zerolog.Logger,
) *Usecase {
	if parentLogger == nil {
		panic("logger cannot be nil")
	}

	logger := parentLogger.With().Str("component", "review_usecase").Logger()

	return &Usecase{
		dictionaryRepo: dictionaryRepo,
		subsRepo:       subsRepo,
		wordStateRepo:  wordStateRepo,
		logger:         &logger,
		sessions:       make(map[int64]*reviewSession),
	}
}

func (u *Usecase) PrepareByDictionaryNumber(
	ctx context.Context,
	userID int64,
	number int,
) (string, error) {
	const op = "PrepareByDictionaryNumber"

	if number <= 0 {
		return "", domain.ErrInvalidDictionaryNumber
	}

	dictionaries, err := u.subsRepo.ListByUser(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if number > len(dictionaries) {
		return "", domain.ErrInvalidDictionaryNumber
	}

	dictionaryID := dictionaries[number-1].ID

	u.clearSession(userID)

	return dictionaryID, nil
}

func (u *Usecase) PrepareByDictionaryID(
	ctx context.Context,
	userID int64,
	dictionaryID string,
) error {
	const op = "PrepareByDictionaryID"

	err := u.prepareByDictionaryIDInner(ctx, userID, dictionaryID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *Usecase) StartForceRound(ctx context.Context, userID int64, dictionaryID string) (*domain.ReviewWord, error) {
	const op = "StartForceRound"

	if err := u.prepareByDictionaryIDInner(ctx, userID, dictionaryID); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now()
	words, err := u.wordStateRepo.ListAllReviewWordsByNearest(ctx, userID, dictionaryID, now)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	// TODO: add transaction for prepareByDictionaryIDInner &
	// ListAllReviewWordsByNearest and get rid of if below:
	// prepareByDictionaryIDInner will check it with HasReviewWords.
	if len(words) == 0 {
		return nil, domain.ErrEmptyReviewWordsList
	}

	u.setSession(userID, dictionaryID, words)

	nextW, err := u.nextWord(userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return nextW, nil
}

func (u *Usecase) prepareByDictionaryIDInner(
	ctx context.Context,
	userID int64,
	dictionaryID string,
) error {
	const op = "prepareByDictionaryIDInner"

	exists, err := u.dictionaryRepo.ExistsByID(ctx, dictionaryID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		return domain.ErrDictionaryNotFound
	}

	subscribed, err := u.subsRepo.IsSubscribedByUser(ctx, userID, dictionaryID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !subscribed {
		return domain.ErrSubscriptionNotFound
	}

	hasReviewWords, err := u.wordStateRepo.HasReviewWords(ctx, userID, dictionaryID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if !hasReviewWords {
		return domain.ErrEmptyReviewWordsList
	}

	u.clearSession(userID)

	return nil
}

func (u *Usecase) StartDueRound(ctx context.Context, userID int64, dictionaryID string) (*domain.ReviewWord /*string,*/, error) {
	const op = "StartDueRound"

	hasReviewWords, err := u.wordStateRepo.HasReviewWords(ctx, userID, dictionaryID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !hasReviewWords {
		return nil, domain.ErrEmptyReviewWordsList
	}

	now := time.Now()
	words, err := u.wordStateRepo.ListDueReviewWords(ctx, userID, dictionaryID, now)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if len(words) == 0 {
		return nil, domain.ErrNoWordsDueForReview
	}

	u.setSession(userID, dictionaryID, words)

	nextW, err := u.nextWord(userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return nextW, nil
}

func (u *Usecase) RateCurrent(ctx context.Context, userID int64, grade int) (*domain.ReviewWord, string, error) {
	const op = "RateCurrent"

	if grade < domain.MinGrade || grade > domain.MaxGrade {
		return nil, "", domain.ErrInvalidReviewGrade
	}

	session, ok := u.getSession(userID)
	if !ok || session.current == nil {
		return nil, "", domain.ErrReviewNotStarted
	}

	now := time.Now()
	result, err := domain.ComputeSM2(&domain.SM2Input{
		EF:           session.current.EF,
		IntervalDays: session.current.IntervalDays,
		Repetition:   session.current.Repetition,
		Grade:        grade,
	}, now)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	if err = u.wordStateRepo.ApplyReviewResult(ctx, &domain.ApplyReviewResultInput{
		UserID:     userID,
		DictWordID: session.current.ID,
		Grade:      grade,
		Result:     result,
		ReviewedAt: now,
	}); err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	nextW, err := u.nextWord(userID)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %w", op, err)
	}

	return nextW, session.dictionaryID, nil
}

func (u *Usecase) Stop(ctx context.Context, userID int64) error {
	const op = "Stop"

	u.clearSession(userID)

	return nil
}


func (u *Usecase) setSession(userID int64, dictionaryID string, words []*domain.ReviewWord) {
	u.sessionMu.Lock()
	defer u.sessionMu.Unlock()

	u.sessions[userID] = &reviewSession{
		dictionaryID: dictionaryID,
		queue:        words,
		current:      nil,
	}
}

func (u *Usecase) getSession(userID int64) (*reviewSession, bool) {
	u.sessionMu.RLock()
	defer u.sessionMu.RUnlock()

	session, ok := u.sessions[userID]
	return session, ok
}

func (u *Usecase) clearSession(userID int64) {
	u.sessionMu.Lock()
	defer u.sessionMu.Unlock()
	delete(u.sessions, userID)
}

func (u *Usecase) nextWord(userID int64) (*domain.ReviewWord, error) {
	u.sessionMu.Lock()
	defer u.sessionMu.Unlock()

	session, ok := u.sessions[userID]
	if !ok || len(session.queue) == 0 {
		delete(u.sessions, userID)
		return nil, domain.ErrReviewRoundFinished
	}

	session.current = session.queue[0]
	session.queue = session.queue[1:]

	return session.current, nil
}
