package service

import (
	"texttool/internal/config"
	"texttool/internal/indexer"
	"texttool/internal/store"
	"texttool/pkg/logger"
)

type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
	index *indexer.Builder
}

func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg, index: indexer.NewBuilder()}
}
