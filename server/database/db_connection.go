package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	Client *mongo.Client
)

func InitializeDB() error {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: unable to find .env file, relying on system environment variables")
	}

	MongoURI := os.Getenv("MONGODB_URI")
	if MongoURI == "" {
		return fmt.Errorf("MONGODB_URI environment variable not set")
	}

	clientOptions := options.Client().ApplyURI(MongoURI)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return fmt.Errorf("Failed to connect to MONGODB: %w", err)
	}

	Client = client
	log.Println("Successfully connected to MongoDB")
	return nil
}

// returns a collection reference from DB
func OpenCollection(collectionName string) (*mongo.Collection, error) {
	if Client == nil {
		return nil, fmt.Errorf("Database client not initialized. Call InitializeDB() first.")
	}

	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		return nil, fmt.Errorf("DATABASE_NAME environment variable not set")
	}

	return Client.Database(databaseName).Collection(collectionName), nil
}
