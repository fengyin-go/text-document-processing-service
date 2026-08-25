package model

import (
	"strings"
	"time"
)

const (
	OpBase64Encode = "base64encode"
	OpBase64Decode = "base64decode"
	OpURLEncode    = "urlencode"
	OpURLDecode    = "urldecode"
	OpUppercase    = "uppercase"
	OpLowercase    = "lowercase"
	OpTrim         = "trim"
	OpReverse      = "reverse"

	StatusPending = "pending"
	StatusDone    = "done"
	StatusFailed  = "failed"
)

var validOperations = map[string]bool{
	OpBase64Encode: true,
	OpBase64Decode: true,
	OpURLEncode:    true,
	OpURLDecode:    true,
	OpUppercase:    true,
	OpLowercase:    true,
	OpTrim:         true,
	OpReverse:      true,
}

var transitions = map[string]map[string]bool{
	StatusPending: {StatusDone: true, StatusFailed: true},
}

func CanTransition(from, to string) bool {
	if m, ok := transitions[from]; ok {
		return m[to]
	}
	return false
}

type TransformTask struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	Operation  string    `json:"operation"`
	Input      string    `json:"input"`
	Output     string    `json:"output"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	DoneAt     *time.Time `json:"done_at,omitempty"`
}

func (t *TransformTask) Validate() error {
	t.Operation = strings.TrimSpace(t.Operation)
	if t.Operation == "" {
		return NewValidationError("operation", "操作类型不能为空")
	}
	if !validOperations[t.Operation] {
		return NewValidationError("operation", "操作类型不合法")
	}
	if t.Status == "" {
		t.Status = StatusPending
	}
	if t.Status != StatusPending && t.Status != StatusDone && t.Status != StatusFailed {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type TransformTaskFilter struct {
	DocumentID string
	Operation  string
	Status     string
}

func (f TransformTaskFilter) Match(t *TransformTask) bool {
	if f.DocumentID != "" && t.DocumentID != f.DocumentID {
		return false
	}
	if f.Operation != "" && t.Operation != f.Operation {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	return true
}
