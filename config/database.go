package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoConfig struct {
	URI      string
	Database string
}

var MongoDBConfig mongoConfig
var MongoDBClient *mongo.Client

func InitMongoDB() {
	// Load configuration
	MongoDBConfig = mongoConfig{
		URI:      GetEnv("MONGODB_URI"),
		Database: GetEnv("MONGODB_DATABASE"),
	}

	fmt.Printf("🔌 Connecting to MongoDB with URI: %s\n", MongoDBConfig.URI)
	clientOptions := options.Client().ApplyURI(MongoDBConfig.URI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("❌ Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ MongoDB is not responding: %v", err)
	}

	fmt.Println("✅ MongoDB connection successful!")
	MongoDBClient = client
}

func GetCollection(collectionName string) *mongo.Collection {
	return MongoDBClient.Database(MongoDBConfig.Database).Collection(collectionName)
}
