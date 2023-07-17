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

func (cs *ChunkService) CreateChunk(hash string) (data.Chunk, error) {
	newChunk := data.Chunk{
		Hash: hash,
	}
	createdChunk, createErr := cs.repo.Add(newChunk)
	if createErr != nil {
		cs.logger.Error("Something went wrong creating chunk", zap.Error(createErr))
		return data.Chunk{}, createErr
	}

	return createdChunk, nil
}
