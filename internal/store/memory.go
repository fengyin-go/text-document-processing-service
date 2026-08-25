package store

import (
	"sync"

	"texttool/internal/model"
)

type MemoryStore struct {
	mu              sync.RWMutex
	textDocuments   map[string]*model.TextDocument
	transformTasks  map[string]*model.TransformTask
	hashRecords     map[string]*model.HashRecord
	regexRules      map[string]*model.RegexRule
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		textDocuments:  make(map[string]*model.TextDocument),
		transformTasks: make(map[string]*model.TransformTask),
		hashRecords:    make(map[string]*model.HashRecord),
		regexRules:     make(map[string]*model.RegexRule),
	}
}

var _ Store = (*MemoryStore)(nil)
