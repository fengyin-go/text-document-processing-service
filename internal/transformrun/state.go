package transformrun

import "sync"

type Attempt struct {
	Version int
	Key     string
}

type State struct {
	mu          sync.Mutex
	version     int
	status      string
	sideEffects int
	keys        map[string]struct{}
}

func NewState() *State {
	return &State{status: "pending", keys: make(map[string]struct{})}
}

func (s *State) Begin(key string) Attempt {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.version++
	s.status = "running"
	if _, exists := s.keys[key]; !exists {
		s.keys[key] = struct{}{}
		s.sideEffects++
	}
	return Attempt{Version: s.version, Key: key}
}

func (s *State) Complete(attempt Attempt, status string) {
	s.mu.Lock()
	s.status = status
	s.mu.Unlock()
}

func (s *State) Snapshot() (string, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status, s.sideEffects
}
