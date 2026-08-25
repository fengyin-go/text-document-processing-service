// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"texttool/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// TextDocument
	CreateTextDocument(d *model.TextDocument) error
	GetTextDocument(id string) (*model.TextDocument, error)
	ListTextDocuments() []*model.TextDocument
	UpdateTextDocument(d *model.TextDocument) error
	DeleteTextDocument(id string) error

	// TransformTask
	CreateTransformTask(t *model.TransformTask) error
	GetTransformTask(id string) (*model.TransformTask, error)
	ListTransformTasks() []*model.TransformTask
	UpdateTransformTask(t *model.TransformTask) error
	DeleteTransformTask(id string) error

	// HashRecord
	CreateHashRecord(h *model.HashRecord) error
	GetHashRecord(id string) (*model.HashRecord, error)
	ListHashRecords() []*model.HashRecord
	UpdateHashRecord(h *model.HashRecord) error
	DeleteHashRecord(id string) error

	// RegexRule
	CreateRegexRule(r *model.RegexRule) error
	GetRegexRule(id string) (*model.RegexRule, error)
	GetRegexRuleByName(name string) (*model.RegexRule, error)
	ListRegexRules() []*model.RegexRule
	UpdateRegexRule(r *model.RegexRule) error
	DeleteRegexRule(id string) error
}
