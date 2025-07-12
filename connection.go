package meerkat

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	defaultClient   *mongo.Client
	connectOnce     sync.Once
	errNotConnected = errors.New("meerkat: not connected to MongoDb, please call meerkat.Connect()")
)

func Connect(ctx context.Context, uri string) error {
	var err error

	connectOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(uri).SetConnectTimeout(10 * time.Second).SetServerSelectionTimeout(30 * time.Second).SetMinPoolSize(5).SetMaxPoolSize((100))

		client, innerErr := mongo.Connect(ctx, clientOptions)

		if innerErr != nil {
			err = innerErr
			return
		}
		if innerErr := client.Ping(ctx, readpref.Primary()); innerErr != nil {
			err = innerErr
			return
		}
		defaultClient = client
	})

	if defaultClient == nil && err == nil {
		return errors.New("meerkat: failed to connect to MongoDb but no specific error was returned")
	}
	return err
}
