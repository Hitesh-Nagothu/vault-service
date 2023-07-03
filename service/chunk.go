package service

import (
	"github.com/Hitesh-Nagothu/vault-service/data"
	"go.uber.org/zap"
)

type ChunkService struct {
	repo   *data.ChunkRepository
	logger *zap.Logger
}

func NewChunkService(logger *zap.Logger, repo *data.ChunkRepository) *ChunkService {
	return &ChunkService{
		logger: logger,
		repo:   repo,
	}
}
