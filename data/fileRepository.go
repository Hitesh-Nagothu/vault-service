package data

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository struct {
	collection *mongo.Collection
}

func NewFileRepository(db *MongoDB) *FileRepository {
	return &FileRepository{
		collection: db.GetDatabase().Collection("file"),
	}
}
