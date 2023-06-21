package data

import (
	"github.com/google/uuid"
)

type FileMetadata struct {
	Name     string
	MimeType string
	Chunks   []uuid.UUID
}

type FileStore struct {
	Data map[uuid.UUID]FileMetadata
}

var storeinstance *FileStore

func GetFileStore() *FileStore {
	if storeinstance == nil {
		storeinstance = &FileStore{
			Data: make(map[uuid.UUID]FileMetadata),
		}
	}
	return storeinstance
}
