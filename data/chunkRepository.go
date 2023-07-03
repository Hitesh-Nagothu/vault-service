package data

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

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
