package service

import (
	"fmt"
	"texttool/internal/transformrun"
)

// RunTransformRetry simulates a retry scenario for a single transform task:
//
//  1. The first attempt is started but its callback is held back (e.g. the
//     downstream call is stuck).
//  2. A retry (second attempt) is started and completes as "done" while the
//     first attempt's callback is still pending.
//  3. The first attempt's stale callback eventually arrives.
//
// The expected outcome is that the retry's terminal status survives the late
// callback, and external processing happens exactly once. Both invariants are
// enforced by State.Complete (see internal/transformrun/state.go): a callback
// from a superseded attempt is ignored, and a terminal status never regresses.
func (s *Service) RunTransformRetry(state *transformrun.State, taskID string, releaseFirst <-chan struct{}) <-chan struct{} {
	first := state.Begin(fmt.Sprintf("%s-attempt-%d", taskID, 1))
	second := state.Begin(fmt.Sprintf("%s-attempt-%d", taskID, 2))
	state.Complete(second, "done")
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-releaseFirst
		// This is the late callback from the first (superseded) attempt. It
		// must be a no-op: status stays "done", sideEffects stays 1.
		state.Complete(first, "running")
	}()
	return done
}
