package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

// Connect initializes the MongoDB connection
func Connect(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	Client = client
	Database = client.Database("dune-form-builder")

	log.Println("✅ Connected to MongoDB successfully!")
	return nil
}

// Disconnect closes the MongoDB connection
func Disconnect() error {
	if Client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := Client.Disconnect(ctx)
	if err != nil {
		return err
	}

	log.Println("📴 Disconnected from MongoDB")
	return nil
}

// Ping checks if the database connection is alive
func Ping() error {
	if Client == nil {
		return mongo.ErrClientDisconnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return Client.Ping(ctx, nil)
}
