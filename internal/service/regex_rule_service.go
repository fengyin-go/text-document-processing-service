package service

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"texttool/internal/model"
	"texttool/pkg/idgen"
)

func (s *Service) CreateRegexRule(input model.RegexRule) (*model.RegexRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := compileRegex(input.Pattern, input.Flags); err != nil {
		return nil, model.NewValidationError("pattern", "正则表达式编译失败: "+err.Error())
	}
	now := time.Now()
	r := &model.RegexRule{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Pattern:     input.Pattern,
		Flags:       input.Flags,
		Description: input.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateRegexRule(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) GetRegexRule(id string) (*model.RegexRule, error) {
	return s.store.GetRegexRule(id)
}

func (s *Service) ListRegexRules(filter model.RegexRuleFilter, page, size int) ([]*model.RegexRule, int, error) {
	all := s.store.ListRegexRules()
	matched := make([]*model.RegexRule, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RegexRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRegexRule(id string, input model.RegexRule) (*model.RegexRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := compileRegex(input.Pattern, input.Flags); err != nil {
		return nil, model.NewValidationError("pattern", "正则表达式编译失败: "+err.Error())
	}
	r, err := s.store.GetRegexRule(id)
	if err != nil {
		return nil, err
	}
	r.Name = input.Name
	r.Pattern = input.Pattern
	r.Flags = input.Flags
	r.Description = input.Description
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateRegexRule(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteRegexRule(id string) error {
	return s.store.DeleteRegexRule(id)
}

// MatchResult 正则匹配结果。
type MatchResult struct {
	Matched bool     `json:"matched"`
	Matches []string `json:"matches"`
}

func (s *Service) MatchWithRule(ruleID, text string) (*MatchResult, error) {
	rule, err := s.store.GetRegexRule(ruleID)
	if err != nil {
		return nil, err
	}
	re, err := compileRegex(rule.Pattern, rule.Flags)
	if err != nil {
		return nil, model.NewValidationError("pattern", "正则表达式编译失败: "+err.Error())
	}
	matches := re.FindAllString(text, -1)
	return &MatchResult{
		Matched: len(matches) > 0,
		Matches: matches,
	}, nil
}

func (s *Service) ReplaceWithRule(ruleID, text, replacement string) (string, error) {
	rule, err := s.store.GetRegexRule(ruleID)
	if err != nil {
		return "", err
	}
	re, err := compileRegex(rule.Pattern, rule.Flags)
	if err != nil {
		return "", model.NewValidationError("pattern", "正则表达式编译失败: "+err.Error())
	}
	return re.ReplaceAllString(text, replacement), nil
}

func compileRegex(pattern, flags string) (*regexp.Regexp, error) {
	var sb strings.Builder
	for _, f := range flags {
		switch f {
		case 'i':
			sb.WriteString("(?i)")
		case 'm':
			sb.WriteString("(?m)")
		case 's':
			sb.WriteString("(?s)")
		}
	}
	sb.WriteString(pattern)
	return regexp.Compile(sb.String())
}

// MatchWithPattern 直接使用模式匹配文本。
func (s *Service) MatchWithPattern(pattern, flags, text string) (*MatchResult, error) {
	re, err := compileRegex(pattern, flags)
	if err != nil {
		return nil, model.NewValidationError("pattern", "正则表达式编译失败: "+err.Error())
	}
	matches := re.FindAllString(text, -1)
	return &MatchResult{
		Matched: len(matches) > 0,
		Matches: matches,
	}, nil
}

// ReplaceWithPattern 直接使用模式替换文本。
func (s *Service) ReplaceWithPattern(pattern, flags, text, replacement string) (string, error) {
	re, err := compileRegex(pattern, flags)
	if err != nil {
		return "", model.NewValidationError("pattern", "正则表达式编译失败: "+err.Error())
	}
	return re.ReplaceAllString(text, replacement), nil
}

// ValidatePattern 验证正则表达式是否合法。
func (s *Service) ValidatePattern(pattern, flags string) error {
	_, err := compileRegex(pattern, flags)
	if err != nil {
		return model.NewValidationError("pattern", "正则表达式编译失败: "+err.Error())
	}
	return nil
}

// ListRuleNames 返回所有规则名称列表。
func (s *Service) ListRuleNames() ([]string, error) {
	all := s.store.ListRegexRules()
	names := make([]string, 0, len(all))
	for _, r := range all {
		names = append(names, r.Name)
	}
	sort.Strings(names)
	return names, nil
}

// GetRuleByName 按名称获取规则。
func (s *Service) GetRuleByName(name string) (*model.RegexRule, error) {
	return s.store.GetRegexRuleByName(name)
}

// BatchMatch 批量匹配多个规则。
func (s *Service) BatchMatch(ruleIDs []string, text string) (map[string]*MatchResult, error) {
	results := make(map[string]*MatchResult)
	for _, id := range ruleIDs {
		result, err := s.MatchWithRule(id, text)
		if err != nil {
			return nil, fmt.Errorf("规则 %s 匹配失败: %w", id, err)
		}
		results[id] = result
	}
	return results, nil
}
