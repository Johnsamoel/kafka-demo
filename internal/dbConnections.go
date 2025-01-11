package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client            *mongo.Client
	ctx               context.Context
	users             *mongo.Collection
	notifications_log *mongo.Collection
	dbName            string
)

const (
	USERS_COLL             = "users"
	NOTIFICATIONS_LOG_COLL = "notifications_log"
)

// ConnectDB establishes a connection to MongoDB and returns the client instance.
func ConnectDB() (*mongo.Client, error) {
	// Load the .env file
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return nil, err
	}

	// Get MongoDB connection string from the .env file
	mongoURI := os.Getenv("MONGO_URI")
	dbName = os.Getenv("DB_NAME")

	if mongoURI == "" || dbName == "" {
		log.Fatalf("Missing required environment variables")
		return nil, fmt.Errorf("missing required environment variables")
	}

	// Set MongoDB client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB Atlas: %v", err)
		return nil, err
	}

	// Ping the MongoDB server to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	fmt.Println("Successfully connected to MongoDB Atlas!")
	return client, nil
}

// GetCollection provides access to a specific MongoDB collection.
func GetCollection(collectionName string) *mongo.Collection {
	if client == nil {
		log.Fatal("Database connection is not established. Call ConnectDB first.")
	}

	db := client.Database(dbName)
	return db.Collection(collectionName)
}

// DisconnectDB disconnects from MongoDB.
func DisconnectDB() {
	if client == nil {
		log.Println("No MongoDB connection to disconnect.")
		return
	}

	// Disconnect the client
	if err := client.Disconnect(ctx); err != nil {
		log.Fatalf("Failed to disconnect MongoDB: %v", err)
	} else {
		fmt.Println("Disconnected from MongoDB Atlas.")
	}
}
