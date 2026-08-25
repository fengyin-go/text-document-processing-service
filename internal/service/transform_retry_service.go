package service

import (
	"fmt"
	"texttool/internal/transformrun"
)

func (s *Service) RunTransformRetry(state *transformrun.State, taskID string, releaseFirst <-chan struct{}) <-chan struct{} {
	first := state.Begin(fmt.Sprintf("%s-attempt-%d", taskID, 1))
	second := state.Begin(fmt.Sprintf("%s-attempt-%d", taskID, 2))
	state.Complete(second, "done")
	done := make(chan struct{})
	go func() {
		defer close(done)
		<-releaseFirst
		state.Complete(first, "running")
	}()
	return done
}
