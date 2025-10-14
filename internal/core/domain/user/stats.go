package user

import (
	"time"
)

type Stats struct {
	Weight    *WeightValue `json:"weight"`
	Height    *HeightValue `json:"height"`
	BFP       *BFP         `json:"bfp"`
	Streak    Streak       `json:"streak"`
	Totals    Totals       `json:"totals"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
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
