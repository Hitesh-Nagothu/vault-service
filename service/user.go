package service

import (
	"github.com/Hitesh-Nagothu/vault-service/data"
	"go.uber.org/zap"
)

type UserService struct {
	repo   *data.UserRepository
	logger *zap.Logger
}

func NewUserService(logger *zap.Logger, repo *data.UserRepository) *UserService {
	return &UserService{
		logger: logger,
		repo:   repo,
	}
}
