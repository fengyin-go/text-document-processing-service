package service

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"texttool/internal/config"
	"texttool/internal/dispatch"
	"texttool/internal/store"
	"texttool/pkg/logger"
)

func TestDispatchCancellationStopsBackendRetries(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{})
	var calls atomic.Int64
	firstCall := make(chan struct{})
	var once sync.Once
	runner := dispatch.NewRunner(8*time.Millisecond, func(context.Context) error {
		calls.Add(1)
		once.Do(func() { close(firstCall) })
		return errors.New("temporary transform backend failure")
	})
	ctx, cancel := context.WithCancel(context.Background())
	done := svc.StartTransformDispatch(ctx, runner)
	<-firstCall
	cancel()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("cancelled transform dispatch did not stop before shutdown: active=%d calls=%d", runner.Active(), calls.Load())
	}
	if calls.Load() != 1 || runner.Active() != 0 {
		t.Fatalf("cancelled transform dispatch kept retrying in the background: active=%d calls=%d", runner.Active(), calls.Load())
	}
}
