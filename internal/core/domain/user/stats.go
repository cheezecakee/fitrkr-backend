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

func NewStats(streak Streak) Stats {
	return Stats{
		Totals:    NewTotals(),
		Streak:    streak,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
