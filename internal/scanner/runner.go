package scanner

import (
	"context"
	"sync"
)

type Backend interface {
	Scan(context.Context, string) (string, error)
}

type BackendFunc func(context.Context, string) (string, error)

func (f BackendFunc) Scan(ctx context.Context, text string) (string, error) {
	return f(ctx, text)
}

type Runner struct {
	mu      sync.Mutex
	backend Backend
	saved   context.Context
}

func NewRunner(backend Backend) *Runner {
	return &Runner{backend: backend}
}

func (r *Runner) Scan(ctx context.Context, text string) (string, error) {
	r.mu.Lock()
	if r.saved == nil {
		r.saved = ctx
	}
	active := r.saved
	r.mu.Unlock()
	if err := active.Err(); err != nil {
		return "", err
	}
	return r.backend.Scan(active, text)
}
