package user

import (
	"time"

	"github.com/google/uuid"
)

type Stats struct {
	ID        uuid.UUID
	Weight    *Weight
	Height    *Height
	BFP       *BFP
	Streak    Streak
	Totals    Totals
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewStats(id uuid.UUID, streak Streak) Stats {
	return Stats{
		ID:        id,
		Totals:    NewTotals(),
		Streak:    streak,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
