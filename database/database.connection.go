package database

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// InitDBConnection creates new mongodb instance
func InitDBConnection(mongoUrl string) *mongo.Client {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")
	return client
}

// OpenCollection opens collection in mongodb
func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("book-DB").Collection(collectionName)
}
