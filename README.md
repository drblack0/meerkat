# Meerkat

Meerkat is an opinionated libaray for working mongoDB, it takes inspiration from mongoose driver for Node js but aims to be lighter and faster than mongoose and still provide all the same functionalities like models, schemas, hooks, validation and connection management. 

For now the meerkat package supports the basic CRUD queries with connection management. 

Here is how to connect with meerkat:

``` package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	// Replace with your actual module path
	"your-app/database" 

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	// 1. Get connection URI from environment variable or use a default
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	
	ctx := context.Background()

	// 2. Connect using the library
	// This returns a singleton instance. Calling it again will return the same instance.
	mongoConn, err := database.Connect(ctx, mongoURI, "myAppDB")
	if err != nil {
		log.Fatalf("Could not initialize database connection: %s\n", err)
	}

	// 3. Defer disconnection for graceful shutdown
	defer func() {
		if err := mongoConn.Disconnect(ctx); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	// 4. Use the connection to get a collection and perform operations
	usersCollection := mongoConn.GetCollection("users")

	// Insert a document
	res, err := usersCollection.InsertOne(ctx, bson.D{
		{Key: "name", Value: "Ada Lovelace"},
		{Key: "email", Value: "ada@example.com"},
	})
	if err != nil {
		log.Fatalf("Failed to insert document: %v", err)
	}

	fmt.Printf("✅ Inserted a single document with ID: %v\n", res.InsertedID)
}
