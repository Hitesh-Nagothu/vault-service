package data

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

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
