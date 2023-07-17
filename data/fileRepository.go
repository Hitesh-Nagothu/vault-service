package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type File struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty"`
	Name     string               `bson:"name"`
	Type     string               `bson:"type"`
	ChunkIDs []primitive.ObjectID `bson:"chunk_ids"`
}

type FileRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewFileRepository(db *MongoDB, logger *zap.Logger) *FileRepository {
	return &FileRepository{
		collection: db.GetDatabase().Collection("file"),
		logger:     logger,
	}
}

func (repo *FileRepository) Add(file File) (File, error) {
	insertResult, err := repo.collection.InsertOne(context.Background(), file)
	if err != nil {
		repo.logger.Fatal("Something went wrong creating the file", zap.Error(err))
		return File{}, nil
	}
	repo.logger.Info("Created a new file successfully", zap.Any("objectId", insertResult.InsertedID))
	insertedFile := File{}
	err = repo.collection.FindOne(context.Background(), bson.M{"_id": insertResult.InsertedID}).Decode(&insertedFile)
	if err != nil {
		repo.logger.Error("Something went wrong getting user by object id", zap.Any("_id", insertResult.InsertedID))
		return File{}, nil
	}
	return insertedFile, nil
}
