package user

import (
	"time"
)

type Stats struct {
	Weight    *WeightValue
	Height    *HeightValue
	BFP       *BFP
	Streak    Streak
	Totals    Totals
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewStats() Stats {
	return Stats{
		Totals:    NewTotals(),
		Streak:    NewStreak(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *Stats) Touch() {
	s.UpdatedAt = time.Now()
}
