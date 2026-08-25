package service

import (
	"context"
	"texttool/internal/dispatch"
)

func (s *Service) StartTransformDispatch(ctx context.Context, runner *dispatch.Runner) <-chan struct{} {
	return runner.Start(context.Background())
}
