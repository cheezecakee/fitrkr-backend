package users

import (
	"context"
	"fmt"
	"time"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
	"github.com/cheezecakee/fitrkr-athena/internal/core/ports"
)

type UpdateBodyMetricsReq struct {
	UserID      string
	WeightValue *float64
	HeightValue *float64
	BFP         *float64
}

func (s *Service) UpdateBodyMetrics(ctx context.Context, req UpdateBodyMetricsReq) error {
	settings, err := s.userRepo.GetSettingsByID(ctx, req.UserID)
	if err != nil {
		logr.Get().Errorf("failed to get user settings: %v", err)
		return fmt.Errorf("failed to get user settings: %w", err)
	}

	var updateBodyMetrics ports.UpdateBodyMetrics

	if req.WeightValue != nil {
		weight, err := user.NewWeight(*req.WeightValue, settings.WeightUnit)
		if err != nil {
			logr.Get().Errorf("failed to create user weight: %v", err)
			return fmt.Errorf("failed to create user weight: %w", err)
		}

		updateBodyMetrics.WeightValue = &weight
	}

	if req.HeightValue != nil {
		height, err := user.NewHeight(*req.HeightValue, settings.HeightUnit)
		if err != nil {
			logr.Get().Errorf("failed to create user height: %v", err)
			return fmt.Errorf("failed to create user height: %w", err)
		}

		updateBodyMetrics.HeightValue = &height
	}

	if req.BFP != nil {
		BFP, err := user.NewBFP(*req.BFP)
		if err != nil {
			logr.Get().Errorf("failed to create user BFP: %v", err)
			return fmt.Errorf("failed to create user BFP: %w", err)
		}

		updateBodyMetrics.BFP = &BFP
	}

	updateBodyMetrics.UpdatedAt = time.Now()

	err = s.userRepo.UpdateBodyMetrics(ctx, updateBodyMetrics, req.UserID)
	if err != nil {
		logr.Get().Errorf("failed to update body metrics %v", err)
		return fmt.Errorf("failed to update body metrics %w", err)
	}

	return nil
}
