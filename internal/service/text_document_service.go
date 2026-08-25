package service

import (
	"sort"
	"strings"
	"time"
	"unicode"

	"texttool/internal/model"
	"texttool/pkg/idgen"
)

func (s *Service) CreateTextDocument(input model.TextDocument) (*model.TextDocument, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	d := &model.TextDocument{
		ID:        idgen.Hex(),
		Title:     input.Title,
		Content:   input.Content,
		Encoding:  input.Encoding,
		SizeBytes: len(input.Content),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateTextDocument(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) GetTextDocument(id string) (*model.TextDocument, error) {
	return s.store.GetTextDocument(id)
}

func (s *Service) ListTextDocuments(filter model.TextDocumentFilter, page, size int) ([]*model.TextDocument, int, error) {
	all := s.store.ListTextDocuments()
	matched := make([]*model.TextDocument, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.TextDocument{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateTextDocument(id string, input model.TextDocument) (*model.TextDocument, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	d, err := s.store.GetTextDocument(id)
	if err != nil {
		return nil, err
	}
	d.Title = input.Title
	d.Content = input.Content
	d.Encoding = input.Encoding
	d.SizeBytes = len(input.Content)
	d.UpdatedAt = time.Now()
	if err := s.store.UpdateTextDocument(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) DeleteTextDocument(id string) error {
	return s.store.DeleteTextDocument(id)
}

// TextStats 返回文本统计信息。
type TextStats struct {
	CharCount  int            `json:"char_count"`
	WordCount  int            `json:"word_count"`
	LineCount  int            `json:"line_count"`
	TopWords   []WordFreq     `json:"top_words"`
}

type WordFreq struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

func (s *Service) AnalyzeText(id string) (*TextStats, error) {
	d, err := s.store.GetTextDocument(id)
	if err != nil {
		return nil, err
	}
	content := d.Content
	stats := &TextStats{
		CharCount: len([]rune(content)),
		LineCount: len(strings.Split(content, "\n")),
	}

	words := extractWords(content)
	stats.WordCount = len(words)

	freq := make(map[string]int)
	for _, w := range words {
		freq[w]++
	}
	topList := make([]WordFreq, 0, len(freq))
	for w, c := range freq {
		topList = append(topList, WordFreq{Word: w, Count: c})
	}
	sort.Slice(topList, func(i, j int) bool {
		if topList[i].Count == topList[j].Count {
			return topList[i].Word < topList[j].Word
		}
		return topList[i].Count > topList[j].Count
	})
	if len(topList) > 10 {
		topList = topList[:10]
	}
	stats.TopWords = topList

	return stats, nil
}

func extractWords(text string) []string {
	var words []string
	var buf strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			buf.WriteRune(r)
		} else if buf.Len() > 0 {
			words = append(words, strings.ToLower(buf.String()))
			buf.Reset()
		}
	}
	if buf.Len() > 0 {
		words = append(words, strings.ToLower(buf.String()))
	}
	return words
}

// DedupLines 按行去重，保留首次出现的顺序。
func (s *Service) DedupLines(id string) (string, error) {
	d, err := s.store.GetTextDocument(id)
	if err != nil {
		return "", err
	}
	lines := strings.Split(d.Content, "\n")
	seen := make(map[string]bool)
	var result []string
	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n"), nil
}
