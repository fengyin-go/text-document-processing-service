package transformrun

import "sync"

type Attempt struct {
	Version int
	Key     string
}

type State struct {
	mu sync.Mutex

	version int
	status  string

	// External processing (e.g. persisting the result) must happen at most once
	// per task, regardless of how many attempts are started or how many late
	// callbacks arrive. sideEffects counts the times it actually fired.
	sideEffects     int
	sideEffectFired bool
}

func NewState() *State {
	return &State{status: "pending"}
}

// Begin starts a new attempt for key, superseding any in-flight attempt.
func (s *State) Begin(key string) Attempt {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.version++
	s.status = "running"
	return Attempt{Version: s.version, Key: key}
}

// Complete applies a callback's result.
//
// Callbacks can arrive out of order: a retry may finish while the attempt it
// supersersed is still stuck, and that stale attempt's callback can land later.
// Two invariants hold here:
//
//  1. A terminal status never regresses. Once a task is done/failed, no late
//     callback may flip it back to running (or to a different terminal).
//  2. External processing fires at most once per task.
//
// Both fall out of ignoring stale/superseded callbacks: only the latest
// attempt is allowed to mutate state, and a terminal status is sticky.
func (s *State) Complete(attempt Attempt, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ignore callbacks from attempts that a newer Begin has superseded.
	// This is what stops a stuck first attempt's late callback from
	// clobbering the retry's terminal result.
	if attempt.Version < s.version {
		return
	}

	// Terminal is sticky for the current attempt too: a duplicate or errant
	// callback must not regress done -> running (or done -> failed). A retry
	// that genuinely needs to change the outcome must Begin a new attempt,
	// which resets status to running.
	if isTerminal(s.status) && status != s.status {
		return
	}

	prev := s.status
	s.status = status

	// Fire external processing exactly once: on the first transition into a
	// terminal state. The guards above ensure a superseded attempt's late
	// completion never reaches here, so it can never double-count.
	if isTerminal(status) && !isTerminal(prev) && !s.sideEffectFired {
		s.sideEffectFired = true
		s.sideEffects++
	}
}

func (s *State) Snapshot() (string, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status, s.sideEffects
}

func isTerminal(status string) bool {
	return status == "done" || status == "failed"
}
