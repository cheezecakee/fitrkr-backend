package users

import (
	"context"
	"fmt"
	"time"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

type GetStatsReq struct {
	ID string
}

type GetStatsResp struct {
	Weight    *user.WeightValue `json:"weight,omitempty"`
	Height    *user.HeightValue `json:"height,omitempty"`
	BFP       *user.BFP         `json:"body_fat_percent,omitempty"`
	Streak    user.Streak       `json:"streak"`
	Total     user.Totals       `json:"total"`
	CreateAt  time.Time         `json:"create_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func (s *Service) GetStats(ctx context.Context, req GetStatsReq) (*GetStatsResp, error) {
	stats, err := s.userRepo.GetStatsByID(ctx, req.ID)
	if err != nil {
		logr.Get().Errorf("failed to get stats: %v", err)
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	settings, err := s.userRepo.GetSettingsByID(ctx, req.ID)
	if err != nil {
		logr.Get().Errorf("failed to get settings: %v", err)
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	resp := &GetStatsResp{
		BFP:       stats.BFP,
		Streak:    stats.Streak,
		Total:     stats.Totals,
		CreateAt:  stats.CreatedAt,
		UpdatedAt: stats.UpdatedAt,
	}

	// Matches user settings metrics

	if stats.Weight != nil {
		displayValue := stats.Weight.Display(settings.WeightUnit)
		resp.Weight = &displayValue
	}

	if stats.Height != nil {
		displayValue := stats.Height.Display(settings.HeightUnit)
		resp.Height = &displayValue
	}

	return resp, nil
}
