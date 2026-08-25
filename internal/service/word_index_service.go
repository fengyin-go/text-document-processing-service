package service

import "texttool/internal/indexer"

func (s *Service) BuildWordIndex(content string, stopWords []string) map[string]int {
	filter := indexer.LoadFilter(stopWords)
	return s.index.Build(content, filter)
}

func (s *Service) LastWordIndex() map[string]int {
	return s.index.Last()
}
