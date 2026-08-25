package service

import (
	"testing"

	"texttool/internal/config"
	"texttool/internal/store"
	"texttool/internal/transformrun"
	"texttool/pkg/logger"
)

func TestTransformRetryRejectsLateAttemptState(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{})
	state := transformrun.NewState()
	releaseFirst := make(chan struct{})
	done := svc.RunTransformRetry(state, "task-42", releaseFirst)
	close(releaseFirst)
	<-done
	status, sideEffects := state.Snapshot()
	if status != "done" || sideEffects != 1 {
		t.Fatalf("late first attempt replaced the retry result or repeated its side effect: status=%s side_effects=%d", status, sideEffects)
	}
}
