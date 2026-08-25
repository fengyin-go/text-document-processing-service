package preview

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrNotReady = errors.New("preview is not ready")

type Artifact struct {
	ID       string
	Title    string
	Segments []string
	render   func() string
}

type Manager struct {
	mu    sync.RWMutex
	cache map[string]*Artifact
}

func NewManager() *Manager {
	return &Manager{cache: make(map[string]*Artifact)}
}

func (m *Manager) Prepare(id, source string) (artifact *Artifact, err error) {
	artifact = &Artifact{ID: id}
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("preview parse failed: %v", recovered)
		}
		m.mu.Lock()
		m.cache[id] = artifact
		m.mu.Unlock()
	}()
	parts := strings.Split(source, "|")
	artifact.Title = parts[0]
	for _, part := range parts[1:] {
		if part == "PANIC" {
			panic("invalid embedded token")
		}
		artifact.Segments = append(artifact.Segments, part)
	}
	artifact.render = func() string {
		return artifact.Title + ":" + strings.Join(artifact.Segments, ",")
	}
	return artifact, nil
}

func (m *Manager) Render(id string) (string, error) {
	m.mu.RLock()
	artifact := m.cache[id]
	m.mu.RUnlock()
	if artifact == nil {
		return "", ErrNotReady
	}
	return artifact.render(), nil
}

func (m *Manager) Delete(id string) {
	m.mu.Lock()
	delete(m.cache, id)
	m.mu.Unlock()
}
