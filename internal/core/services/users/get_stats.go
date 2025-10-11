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
	Weight           *user.Weight `json:"weight,omitempty"`
	Height           *user.Height `json:"height,omitempty"`
	BFP              *float64     `json:"body_fat_percent,omitempty"`
	RestDays         int          `json:"rest_days"`
	CurrentStreak    int          `json:"current_streak"`
	LongestStreak    int          `json:"longest_streak"`
	TotalWorkouts    int          `json:"total_workouts"`
	TotalLifted      float64      `json:"total_lifted"`
	TotalTimeMinutes int          `json:"total_time_minutes"`
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
		BFP:              stats.BFP,
		RestDays:         stats.RestDays,
		CurrentStreak:    stats.Current,
		LongestStreak:    stats.Longest,
		TotalWorkouts:    stats.TotalWorkouts,
		TotalLifted:      stats.TotalLifted,
		TotalTimeMinutes: stats.TotalTime,
	}

	// Matches user settings metrics

	if stats.WeightValue != nil {
		w := user.Weight{Value: *stats.WeightValue}
		resp.Weight = &user.Weight{Value: w.Display(settings.WeightUnit)}
	}

	if stats.HeightValue != nil {
		h := &user.Height{Value: *stats.HeightValue}
		resp.Height = &user.Height{Value: h.Display(settings.HeightUnit)}
	}

	return resp, nil
}
