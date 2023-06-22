package data

import (
	"github.com/google/uuid"
)

type ChunkStore struct {
	Data map[uuid.UUID]string
}

var chunkStoreInstance *ChunkStore

func GetChunkStore() *ChunkStore {
	if chunkStoreInstance == nil {
		chunkStoreInstance = &ChunkStore{
			Data: make(map[uuid.UUID]string),
		}
	}

	return chunkStoreInstance
}
