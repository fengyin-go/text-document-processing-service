package service

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"time"

	"texttool/internal/model"
	"texttool/pkg/idgen"
)

func (s *Service) CreateHashRecord(input model.HashRecord) (*model.HashRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTextDocument(input.DocumentID); err != nil {
		return nil, model.NewValidationError("document_id", "关联文档不存在")
	}
	now := time.Now()
	h := &model.HashRecord{
		ID:         idgen.Hex(),
		DocumentID: input.DocumentID,
		Algorithm:  input.Algorithm,
		Hash:       input.Hash,
		CreatedAt:  now,
	}
	if err := s.store.CreateHashRecord(h); err != nil {
		return nil, err
	}
	return h, nil
}

func (s *Service) ComputeHash(documentID, algorithm string) (*model.HashRecord, error) {
	if algorithm != model.AlgoMD5 && algorithm != model.AlgoSHA1 && algorithm != model.AlgoSHA256 {
		return nil, model.NewValidationError("algorithm", "算法不合法")
	}
	d, err := s.store.GetTextDocument(documentID)
	if err != nil {
		return nil, err
	}
	var hashStr string
	switch algorithm {
	case model.AlgoMD5:
		sum := md5.Sum([]byte(d.Content))
		hashStr = hex.EncodeToString(sum[:])
	case model.AlgoSHA1:
		sum := sha1.Sum([]byte(d.Content))
		hashStr = hex.EncodeToString(sum[:])
	case model.AlgoSHA256:
		sum := sha256.Sum256([]byte(d.Content))
		hashStr = hex.EncodeToString(sum[:])
	}
	now := time.Now()
	h := &model.HashRecord{
		ID:         idgen.Hex(),
		DocumentID: documentID,
		Algorithm:  algorithm,
		Hash:       hashStr,
		CreatedAt:  now,
	}
	if err := s.store.CreateHashRecord(h); err != nil {
		return nil, err
	}
	return h, nil
}

func (s *Service) GetHashRecord(id string) (*model.HashRecord, error) {
	return s.store.GetHashRecord(id)
}

func (s *Service) ListHashRecords(filter model.HashRecordFilter, page, size int) ([]*model.HashRecord, int, error) {
	all := s.store.ListHashRecords()
	matched := make([]*model.HashRecord, 0, len(all))
	for _, h := range all {
		if filter.Match(h) {
			matched = append(matched, h)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.HashRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteHashRecord(id string) error {
	return s.store.DeleteHashRecord(id)
}

func (s *Service) VerifyHash(documentID, algorithm, expected string) (bool, error) {
	h, err := s.ComputeHash(documentID, algorithm)
	if err != nil {
		return false, err
	}
	return h.Hash == expected, nil
}
