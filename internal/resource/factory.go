package resource

import (
	"errors"
	"sync"
)

var (
	ErrLimit    = errors.New("text source limit reached")
	ErrNotFound = errors.New("text source not found")
)

type Factory struct {
	mu      sync.Mutex
	maxOpen int
	open    int
	data    map[string]string
}

func NewFactory(maxOpen int, data map[string]string) *Factory {
	return &Factory{maxOpen: maxOpen, data: data}
}

func (f *Factory) Open(name string) (*Handle, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	content, ok := f.data[name]
	if !ok {
		return nil, ErrNotFound
	}
	if f.open >= f.maxOpen {
		return nil, ErrLimit
	}
	f.open++
	return &Handle{factory: f, content: content}, nil
}

func (f *Factory) OpenCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.open
}

type Handle struct {
	factory *Factory
	content string
	closed  bool
}

func (h *Handle) Read() string { return h.content }

func (h *Handle) Close() {
	h.factory.mu.Lock()
	defer h.factory.mu.Unlock()
	if h.closed {
		return
	}
	h.closed = true
	h.factory.open--
}
