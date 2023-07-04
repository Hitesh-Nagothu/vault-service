package service

import (
	"time"

	"github.com/Hitesh-Nagothu/vault-service/data"
	"github.com/google/uuid"
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

func (service *UserService) CreateUser(email string) (data.User, error) {
	//check if user with same email already exists
	_, err := service.getUser(email)
	if err != nil {
		service.logger.Error("Something went wrong creating the user", zap.Error(err))
		return data.User{}, nil
	}

	newUser := data.User{
		Email:          email,
		LastAccessedOn: time.Now(),
		Files:          []uuid.UUID{},
	}
	createdUser, createErr := service.repo.Add(&newUser)
	if createErr != nil {
		service.logger.Error("Failed to create new user", zap.Error(createErr))
		return data.User{}, nil
	}
	return createdUser, nil
}

func (service *UserService) getUser(email string) (data.User, error) {
	//preprocess email
	user, err := service.repo.Get(email)
	if err != nil {
		return data.User{}, err
	}

	return user, nil
}
