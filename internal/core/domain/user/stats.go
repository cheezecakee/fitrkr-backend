package user

import (
	"time"
)

type Stats struct {
	Weight    *Weight
	Height    *Height
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
