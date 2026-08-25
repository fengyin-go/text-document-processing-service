package service

import (
	"context"
	"texttool/internal/batchmatch"
)

func (s *Service) MatchRuleBatch(ctx context.Context, coordinator *batchmatch.Coordinator, ids []string) ([]string, error) {
	result, err := coordinator.Run(ctx, ids)
	if err != nil && ctx.Err() == nil {
		return coordinator.Run(ctx, ids)
	}
	return result, err
}
