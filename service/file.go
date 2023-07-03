package service

import (
	"github.com/Hitesh-Nagothu/vault-service/data"
	"go.uber.org/zap"
)

type FileService struct {
	repo   *data.FileRepository
	logger *zap.Logger
}

func NewFileService(logger *zap.Logger, repo *data.FileRepository) *FileService {
	return &FileService{
		logger: logger,
		repo:   repo,
	}
}
