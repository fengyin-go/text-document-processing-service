package store

import (
	"texttool/internal/model"
)

func (s *MemoryStore) CreateTextDocument(d *model.TextDocument) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.textDocuments[d.ID] = d
	return nil
}

func (s *MemoryStore) GetTextDocument(id string) (*model.TextDocument, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.textDocuments[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

func (s *MemoryStore) ListTextDocuments() []*model.TextDocument {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TextDocument, 0, len(s.textDocuments))
	for _, d := range s.textDocuments {
		list = append(list, d)
	}
	return list
}

func (s *MemoryStore) UpdateTextDocument(d *model.TextDocument) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.textDocuments[d.ID]; !ok {
		return ErrNotFound
	}
	s.textDocuments[d.ID] = d
	return nil
}

func (s *MemoryStore) DeleteTextDocument(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.textDocuments[id]; !ok {
		return ErrNotFound
	}
	delete(s.textDocuments, id)
	return nil
}
