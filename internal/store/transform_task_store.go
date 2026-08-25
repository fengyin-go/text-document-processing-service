package store

import (
	"texttool/internal/model"
)

func (s *MemoryStore) CreateTransformTask(t *model.TransformTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.transformTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) GetTransformTask(id string) (*model.TransformTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.transformTasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListTransformTasks() []*model.TransformTask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TransformTask, 0, len(s.transformTasks))
	for _, t := range s.transformTasks {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateTransformTask(t *model.TransformTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.transformTasks[t.ID]; !ok {
		return ErrNotFound
	}
	s.transformTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteTransformTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.transformTasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.transformTasks, id)
	return nil
}
