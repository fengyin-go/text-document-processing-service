package service

import (
	"texttool/internal/config"
	"texttool/internal/preview"
	"texttool/internal/store"
	"texttool/pkg/logger"
)

type Service struct {
	store   store.Store
	log     *logger.Logger
	cfg     *config.Config
	preview *preview.Manager
}

func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg, preview: preview.NewManager()}
}
