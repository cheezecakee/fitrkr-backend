package users

import (
	"context"
	"fmt"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

type GetSettingsReq struct {
	ID string
}

type GetSettingsResp struct {
	Settings user.Settings
}

func (s *Service) GetSettings(ctx context.Context, req GetSettingsReq) (*GetSettingsResp, error) {
	settings, err := s.userRepo.GetSettingsByID(ctx, req.ID)
	if err != nil {
		logr.Get().Errorf("failed to get settings: %v", err)
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	return &GetSettingsResp{Settings: *settings}, nil
}
