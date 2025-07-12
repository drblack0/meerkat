package meerkat

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoConnection struct {
	Client *mongo.Client
	Db     *mongo.Database
}

// singleton instance
var instance *MongoConnection
var once sync.Once

func Connect(ctx context.Context, uri, dbName string) (*MongoConnection, error) {
	var connectErr error

	// once.Do will ensure the connection logic is executed only once
	// across all goroutines.
	once.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(uri)

		// Set a timeout for the connection attempt.
		connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		// Connect to MongoDB
		client, err := mongo.Connect(connectCtx, clientOptions)
		if err != nil {
			connectErr = fmt.Errorf("meerkat: failed to connect to mongo: %w", err)
			return
		}

		// Ping the primary to verify the connection.
		// This is a crucial step to ensure the connection is truly established.
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
			connectErr = fmt.Errorf("meerkat: failed to ping mongo: %w", err)
			// Disconnect if ping fails to clean up resources
			_ = client.Disconnect(context.Background())
			return
		}

		fmt.Println("meerkat: Successfully connected to MongoDB!")

		// Create the singleton instance
		instance = &MongoConnection{
			Client: client,
			Db:     client.Database(dbName),
		}
	})

	if connectErr != nil {
		return nil, connectErr
	}

	if instance == nil {
		// This can happen if another goroutine called Connect, but it failed.
		// We need to reset `once` to allow for another connection attempt.
		// Note: This is an advanced use case for handling retries at the application level.
		once = sync.Once{}
		return nil, fmt.Errorf("meerkat: failed to get mongo instance, connection might have failed previously")
	}

	return instance, nil
}

func (mc *MongoConnection) GetCollection(name string) *mongo.Collection {
	return mc.Db.Collection(name)
}

func (mc *MongoConnection) GetClient() *mongo.Client {
	return mc.Client
}
