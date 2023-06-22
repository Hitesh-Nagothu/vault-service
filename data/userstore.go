package data

import "github.com/google/uuid"

type UserStore struct {
	Data map[string]UserMetadata
}

type UserMetadata struct {
	LastAccessTime string
	Files          []uuid.UUID
}

var userStoreInstance *UserStore

func GetUserStore() *UserStore {
	if userStoreInstance == nil {
		userStoreInstance = &UserStore{
			Data: make(map[string]UserMetadata),
		}
	}
	return userStoreInstance
}
