package service

import (
	"encoding/base64"
	"net/url"
	"sort"
	"strings"
	"time"

	"texttool/internal/model"
	"texttool/pkg/idgen"
)

func (s *Service) CreateTransformTask(input model.TransformTask) (*model.TransformTask, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if input.DocumentID != "" {
		if _, err := s.store.GetTextDocument(input.DocumentID); err != nil {
			return nil, model.NewValidationError("document_id", "关联文档不存在")
		}
	}
	now := time.Now()
	t := &model.TransformTask{
		ID:         idgen.Hex(),
		DocumentID: input.DocumentID,
		Operation:  input.Operation,
		Input:      input.Input,
		Status:     model.StatusPending,
		CreatedAt:  now,
	}
	if err := s.store.CreateTransformTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetTransformTask(id string) (*model.TransformTask, error) {
	return s.store.GetTransformTask(id)
}

func (s *Service) ListTransformTasks(filter model.TransformTaskFilter, page, size int) ([]*model.TransformTask, int, error) {
	all := s.store.ListTransformTasks()
	matched := make([]*model.TransformTask, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.TransformTask{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) ExecuteTransformTask(id string) (*model.TransformTask, error) {
	t, err := s.store.GetTransformTask(id)
	if err != nil {
		return nil, err
	}
	if t.Status != model.StatusPending {
		return nil, model.NewValidationError("status", "任务不处于待处理状态")
	}

	output, err := applyOperation(t.Operation, t.Input)
	if err != nil {
		t.Status = model.StatusFailed
		t.Output = err.Error()
	} else {
		t.Status = model.StatusDone
		t.Output = output
		now := time.Now()
		t.DoneAt = &now
	}
	if err := s.store.UpdateTransformTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func applyOperation(op, input string) (string, error) {
	switch op {
	case model.OpBase64Encode:
		return base64.StdEncoding.EncodeToString([]byte(input)), nil
	case model.OpBase64Decode:
		b, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case model.OpURLEncode:
		return url.QueryEscape(input), nil
	case model.OpURLDecode:
		return url.QueryUnescape(input)
	case model.OpUppercase:
		return strings.ToUpper(input), nil
	case model.OpLowercase:
		return strings.ToLower(input), nil
	case model.OpTrim:
		return strings.TrimSpace(input), nil
	case model.OpReverse:
		runes := []rune(input)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes), nil
	default:
		return "", model.NewValidationError("operation", "不支持的操作类型")
	}
}

func (s *Service) DeleteTransformTask(id string) error {
	return s.store.DeleteTransformTask(id)
}
