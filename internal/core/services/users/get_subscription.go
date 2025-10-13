package users

import (
	"context"
	"fmt"

	"github.com/cheezecakee/logr"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

type GetSubscriptionReq struct {
	ID string
}

type GetSubscriptionResp struct {
	Subscription user.Subscription
}

func (s *Service) GetSubscription(ctx context.Context, req GetSubscriptionReq) (*GetSubscriptionResp, error) {
	sub, err := s.userRepo.GetSubscriptionByID(ctx, req.ID)
	if err != nil {
		logr.Get().Errorf("failed to get subscription: %v", err)
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &GetSubscriptionResp{Subscription: *sub}, nil
}
