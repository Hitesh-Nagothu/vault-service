package data

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *MongoDB) *UserRepository {
	return &UserRepository{
		collection: db.GetDatabase().Collection("user"),
	}
}
