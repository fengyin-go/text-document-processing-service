package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"texttool/internal/batchmatch"
	"texttool/internal/config"
	"texttool/internal/store"
	"texttool/pkg/logger"
)

func TestBatchMatchStopsOnResolverFailure(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{})
	missing := errors.New("rule missing")
	coordinator := batchmatch.NewCoordinator(batchmatch.ResolverFunc(func(id string) (string, error) {
		if id == "missing" {
			return "", missing
		}
		return strings.ToUpper(id), nil
	}))
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	_, err := svc.MatchRuleBatch(ctx, coordinator, []string{"alpha", "missing", "omega"})
	if !errors.Is(err, missing) {
		t.Fatalf("resolver failure turned into a batch timeout instead of returning the rule error: err=%v active=%d", err, coordinator.Active())
	}
	deadline := time.Now().Add(50 * time.Millisecond)
	for coordinator.Active() != 0 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if coordinator.Active() != 0 {
		t.Fatalf("resolver failure left batch producer or collector running: active=%d", coordinator.Active())
	}
}
