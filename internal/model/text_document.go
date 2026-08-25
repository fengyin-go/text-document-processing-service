package model

import (
	"strings"
	"time"
)

const (
	EncodingUTF8  = "utf8"
	EncodingASCII = "ascii"
)

type TextDocument struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Encoding  string    `json:"encoding"`
	SizeBytes int       `json:"size_bytes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *TextDocument) Validate() error {
	d.Title = strings.TrimSpace(d.Title)
	d.Content = strings.TrimSpace(d.Content)
	d.Encoding = strings.TrimSpace(d.Encoding)
	if d.Title == "" {
		return NewValidationError("title", "标题不能为空")
	}
	if d.Encoding == "" {
		d.Encoding = EncodingUTF8
	}
	if d.Encoding != EncodingUTF8 && d.Encoding != EncodingASCII {
		return NewValidationError("encoding", "编码类型不合法")
	}
	if d.SizeBytes < 0 {
		return NewValidationError("size_bytes", "大小不能为负数")
	}
	return nil
}

type TextDocumentFilter struct {
	Encoding string
	Keyword  string
}

func (f TextDocumentFilter) Match(d *TextDocument) bool {
	if f.Encoding != "" && d.Encoding != f.Encoding {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(d.Title), k) {
			return false
		}
	}
	return true
}
