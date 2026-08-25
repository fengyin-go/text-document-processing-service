package model

import (
	"strings"
	"time"
)

type RegexRule struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Pattern     string    `json:"pattern"`
	Flags       string    `json:"flags"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (r *RegexRule) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Pattern = strings.TrimSpace(r.Pattern)
	if r.Name == "" {
		return NewValidationError("name", "规则名称不能为空")
	}
	if r.Pattern == "" {
		return NewValidationError("pattern", "正则表达式不能为空")
	}
	return nil
}

type RegexRuleFilter struct {
	Keyword string
}

func (f RegexRuleFilter) Match(r *RegexRule) bool {
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Name), k) &&
			!strings.Contains(strings.ToLower(r.Description), k) {
			return false
		}
	}
	return true
}
