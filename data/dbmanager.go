package data

import (
	"context"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoDBInstance *MongoDB
	mongoDBOnce     sync.Once
)

type Database interface {
	GetConnection() *mongo.Client
}

type MongoDB struct {
	connection *mongo.Client
}

func GetMongoDBInstance() *MongoDB {
	mongoDBOnce.Do(func() {
		mongoDBInstance = createMongoDBConnection()
	})
	return mongoDBInstance
}

func createMongoDBConnection() *MongoDB {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	return &MongoDB{
		connection: client,
	}
}

func (db *MongoDB) GetDatabase() *mongo.Database {
	return db.connection.Database("vault")
}
