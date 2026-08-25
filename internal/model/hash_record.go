package model

import (
	"strings"
	"time"
)

const (
	AlgoMD5    = "md5"
	AlgoSHA1   = "sha1"
	AlgoSHA256 = "sha256"
)

type HashRecord struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	Algorithm  string    `json:"algorithm"`
	Hash       string    `json:"hash"`
	CreatedAt  time.Time `json:"created_at"`
}

func (h *HashRecord) Validate() error {
	h.Algorithm = strings.TrimSpace(h.Algorithm)
	if h.Algorithm == "" {
		return NewValidationError("algorithm", "算法不能为空")
	}
	if h.Algorithm != AlgoMD5 && h.Algorithm != AlgoSHA1 && h.Algorithm != AlgoSHA256 {
		return NewValidationError("algorithm", "算法不合法")
	}
	if h.DocumentID == "" {
		return NewValidationError("document_id", "文档 ID 不能为空")
	}
	return nil
}

type HashRecordFilter struct {
	DocumentID string
	Algorithm  string
}

func (f HashRecordFilter) Match(h *HashRecord) bool {
	if f.DocumentID != "" && h.DocumentID != f.DocumentID {
		return false
	}
	if f.Algorithm != "" && h.Algorithm != f.Algorithm {
		return false
	}
	return true
}
