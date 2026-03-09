package domain

import (
	"fmt"
	"math"
	"time"
)

const (
	MinGrade = 0
	MaxGrade = 3
)

type SM2Input struct {
	EF           float64
	IntervalDays int
	Repetition   int
	Grade        int
}

type SM2Result struct {
	EF           float64
	IntervalDays int
	Repetition   int
	NextReviewAt time.Time
}

type ApplyReviewResultInput struct {
	UserID     int64
	DictWordID string
	Grade      int
	Result     *SM2Result
	ReviewedAt time.Time
}

func ComputeSM2(input *SM2Input, now time.Time) (*SM2Result, error) {
	if input == nil {
		return nil, fmt.Errorf("compute sm2: input is nil")
	}

	if input.Grade < MinGrade || input.Grade > MaxGrade {
		return nil, ErrInvalidReviewGrade
	}

	ef := input.EF
	if ef < 1.3 {
		ef = 2.5
	}

	repetition := input.Repetition
	interval := input.IntervalDays
	if interval < 0 {
		interval = 0
	}

	if input.Grade < 3 {
		repetition = 0
		interval = 1
	} else {
		switch repetition {
		case 0:
			interval = 1
		case 1:
			interval = 6
		default:
			interval = int(math.Round(float64(interval) * ef))
			if interval < 1 {
				interval = 1
			}
		}

		repetition++
	}

	qualityDiff := float64(MaxGrade - input.Grade)
	ef += (0.1 - qualityDiff*(0.08+qualityDiff*0.02))
	if ef < 1.3 {
		ef = 1.3
	}

	return &SM2Result{
		EF:           ef,
		IntervalDays: interval,
		Repetition:   repetition,
		NextReviewAt: now.AddDate(0, 0, interval),
	}, nil
}
