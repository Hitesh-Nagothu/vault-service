package data

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type ChunkRepository struct {
	collection *mongo.Collection
}

func NewChunkRepository(db *MongoDB) *ChunkRepository {
	return &ChunkRepository{
		collection: db.GetDatabase().Collection("chunk"),
	}
}
