package users

import (
	"context"
	"fmt"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

type GetStatsReq struct {
	ID string
}

type GetStatsResp struct {
	Stats user.Stats
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
		Stats: *stats,
	}

	// Matches user settings metrics

	if stats.Weight != nil {
		displayValue := stats.Weight.Display(settings.WeightUnit)
		resp.Stats.Weight = &displayValue
	}

	if stats.Height != nil {
		displayValue := stats.Height.Display(settings.HeightUnit)
		resp.Stats.Height = &displayValue
	}

	return resp, nil
}
