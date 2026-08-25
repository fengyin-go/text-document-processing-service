package store

import (
	"strings"

	"texttool/internal/model"
)

func (s *MemoryStore) CreateRegexRule(r *model.RegexRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.regexRules {
		if strings.EqualFold(exist.Name, r.Name) {
			return ErrConflict
		}
	}
	s.regexRules[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRegexRule(id string) (*model.RegexRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.regexRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) GetRegexRuleByName(name string) (*model.RegexRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.regexRules {
		if strings.EqualFold(r.Name, name) {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListRegexRules() []*model.RegexRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RegexRule, 0, len(s.regexRules))
	for _, r := range s.regexRules {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateRegexRule(r *model.RegexRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.regexRules[r.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.regexRules {
		if exist.ID != r.ID && strings.EqualFold(exist.Name, r.Name) {
			return ErrConflict
		}
	}
	s.regexRules[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRegexRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.regexRules[id]; !ok {
		return ErrNotFound
	}
	delete(s.regexRules, id)
	return nil
}
