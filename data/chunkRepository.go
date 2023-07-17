package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Chunk struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Hash string             `bson:"hash"`
}
type ChunkRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewChunkRepository(db *MongoDB, logger *zap.Logger) *ChunkRepository {
	return &ChunkRepository{
		collection: db.GetDatabase().Collection("chunk"),
		logger:     logger,
	}
}

// returns the object id of the inserted chunk hash
func (repo *ChunkRepository) Add(chunk Chunk) (Chunk, error) {
	insertResult, err := repo.collection.InsertOne(context.Background(), chunk)
	if err != nil {
		repo.logger.Fatal("Something went wrong creating the user", zap.Error(err))
		return Chunk{}, nil
	}
	repo.logger.Info("Created a new chunk successfully", zap.Any("objectId", insertResult.InsertedID))
	insertedChunk := Chunk{}
	err = repo.collection.FindOne(context.Background(), bson.M{"_id": insertResult.InsertedID}).Decode(&insertedChunk)
	if err != nil {
		repo.logger.Error("Something went wrong getting user by object id", zap.Any("_id", insertResult.InsertedID))
		return Chunk{}, nil
	}
	return insertedChunk, nil
}
