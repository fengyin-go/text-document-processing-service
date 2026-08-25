package service

import (
	"context"
	"errors"
	"texttool/internal/batchmatch"
)

func (s *Service) MatchRuleBatch(ctx context.Context, coordinator *batchmatch.Coordinator, ids []string) ([]string, error) {
	// 解析失败（如规则不存在）应立即返回该错误，不得重试。
	// 仅当上下文被取消时才按原语义重跑一次。
	result, err := coordinator.Run(ctx, ids)
	if err != nil && errors.Is(err, context.Canceled) && ctx.Err() != nil {
		return coordinator.Run(ctx, ids)
	}
	return result, err
}
