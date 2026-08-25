package store

import (
	"texttool/internal/model"
)

func (s *MemoryStore) CreateHashRecord(h *model.HashRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hashRecords[h.ID] = h
	return nil
}

func (s *MemoryStore) GetHashRecord(id string) (*model.HashRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h, ok := s.hashRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return h, nil
}

func (s *MemoryStore) ListHashRecords() []*model.HashRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.HashRecord, 0, len(s.hashRecords))
	for _, h := range s.hashRecords {
		list = append(list, h)
	}
	return list
}

func (s *MemoryStore) UpdateHashRecord(h *model.HashRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.hashRecords[h.ID]; !ok {
		return ErrNotFound
	}
	s.hashRecords[h.ID] = h
	return nil
}

func (s *MemoryStore) DeleteHashRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.hashRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.hashRecords, id)
	return nil
}
