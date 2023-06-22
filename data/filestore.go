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

var fileStoreInstance *FileStore

func GetFileStore() *FileStore {
	if fileStoreInstance == nil {
		fileStoreInstance = &FileStore{
			Data: make(map[uuid.UUID]FileMetadata),
		}
	}
	return fileStoreInstance
}
