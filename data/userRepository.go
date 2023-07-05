package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type User struct {
	Email          string               `bson:"email"`
	LastAccessedOn time.Time            `bson:"last_accessed_on"`
	Files          []primitive.ObjectID `bson:"files"`
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
	insertResult, err := repo.collection.InsertOne(context.Background(), user)
	if err != nil {
		repo.logger.Fatal("Something went wrong creating the user", zap.Error(err))
		return User{}, nil
	}
	repo.logger.Info("Created a new user successfully", zap.Any("objectId", insertResult.InsertedID))
	insertedUser := User{}
	err = repo.collection.FindOne(context.Background(), bson.M{"_id": insertResult.InsertedID}).Decode(&insertedUser)
	if err != nil {
		repo.logger.Error("Something went wrong getting user by object id", zap.Any("_id", insertResult.InsertedID))
		return User{}, nil
	}
	return insertedUser, nil
}

func (repo *UserRepository) Get(email string) (User, error) {

	// Define the filter to match the email field
	filter := bson.M{"email": email}

	// Query the collection for a matching document
	var result User
	err := repo.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			repo.logger.Error("Email ID not found", zap.String("email", email))
			return User{}, nil //not propagating the error to the caller, to allow creating new if not found
		} else {
			repo.logger.Fatal("Something went wrong getting user with email", zap.String("email", email), zap.Error(err))
			return User{}, err
		}
	}

	return result, nil
}
