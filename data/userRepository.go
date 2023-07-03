package data

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type UserRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewUserRepository(db *MongoDB, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		collection: db.GetDatabase().Collection("user"),
		logger:     logger,
	}
}
