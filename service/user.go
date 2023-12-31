package service

import (
	"errors"
	"time"

	"github.com/Hitesh-Nagothu/vault-service/data"
	"github.com/Hitesh-Nagothu/vault-service/utility"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	existingUser, err := service.GetUser(email)
	if err != nil {
		service.logger.Error("Something went wrong creating the user", zap.Error(err))
		return data.User{}, nil
	}

	isEmptyUserData := utility.IsStructEmpty(existingUser)
	if !isEmptyUserData {
		service.logger.Error("User with email already exists", zap.String("email", email))
		return data.User{}, errors.New("user with email already exists")
	}

	newUser := data.User{
		Email:          email,
		LastAccessedOn: time.Now(),
		Files:          []primitive.ObjectID{},
	}
	createdUser, createErr := service.repo.Add(&newUser)
	if createErr != nil {
		service.logger.Error("Failed to create new user", zap.Error(createErr))
		return data.User{}, nil
	}
	return createdUser, nil
}

func (service *UserService) GetUser(email string) (data.User, error) {

	user, err := service.repo.Get(email)
	if err != nil {
		return data.User{}, err
	}

	return user, nil
}

func (service *UserService) UpdateUser(userId primitive.ObjectID, userUpdateBody data.User) error {
	return service.repo.Update(userId, userUpdateBody)
}
