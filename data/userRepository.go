package data

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type User struct {
	Email          string      `bson:"email"`
	LastAccessedOn time.Time   `bson:"last_accessed_on"`
	Files          []uuid.UUID `bson:"files"`
}

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

func (repo *UserRepository) Add(user *User) (User, error) {
	// Insert the user document into the users collection
	_, err := repo.collection.InsertOne(context.Background(), user)
	if err != nil {
		repo.logger.Fatal("Somethign went wrong saving the user", zap.Error(err))
		return User{}, nil
	}
	return User{}, nil
}

func (repo *UserRepository) Get(email string) (User, error) {

	// Define the filter to match the email field
	filter := bson.M{"email": email}

	// Query the collection for a matching document
	var result User
	err := repo.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			repo.logger.Error("Email ID not found")
		} else {
			repo.logger.Fatal("Something went wrong", zap.Error(err))
		}
		return User{}, err
	}

	return result, nil
}
