package data

import (
	"github.com/google/uuid"
)

type ChunkStore struct {
	Data map[uuid.UUID]string
}

var storeInstance *ChunkStore

func GetChunkStore() *ChunkStore {
	if storeInstance == nil {
		storeInstance = &ChunkStore{
			Data: make(map[uuid.UUID]string),
		}
	}

	return storeInstance
}
