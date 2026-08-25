package indexer

import (
	"strings"
	"sync"
)

type Filter interface {
	Allow(string) bool
}

type stopFilter struct {
	blocked map[string]struct{}
}

func (f *stopFilter) Allow(word string) bool {
	_, blocked := f.blocked[strings.ToLower(word)]
	return !blocked
}

func LoadFilter(words []string) Filter {
	var filter *stopFilter
	if len(words) == 0 {
		return filter
	}
	filter = &stopFilter{blocked: make(map[string]struct{}, len(words))}
	for _, word := range words {
		filter.blocked[strings.ToLower(word)] = struct{}{}
	}
	return filter
}

type Builder struct {
	mu   sync.RWMutex
	last map[string]int
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Build(content string, filter Filter) map[string]int {
	result := make(map[string]int)
	b.mu.Lock()
	b.last = result
	b.mu.Unlock()
	for _, word := range strings.Fields(content) {
		result[strings.ToLower(word)]++
		if !filter.Allow(word) {
			delete(result, strings.ToLower(word))
		}
	}
	return result
}

func (b *Builder) Last() map[string]int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	result := make(map[string]int, len(b.last))
	for word, count := range b.last {
		result[word] = count
	}
	return result
}
