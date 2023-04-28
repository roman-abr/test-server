package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
	return client
}

var client *mongo.Client = ConnectDB()

func Collection(collectionName string) *mongo.Collection {
	collection := client.Database("test-server").Collection(collectionName)
	return collection
}
