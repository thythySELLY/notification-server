package config

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoConfig struct {
	URI      string
	Database string
}

// Biến toàn cục cho MongoDB
var (
	MongoDBConfig mongoConfig
	MongoDBClient *mongo.Client
	once          sync.Once
)

func InitMongoDB() {
	once.Do(func() { 
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

		// Kiểm tra kết nối MongoDB
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("❌ MongoDB is not responding: %v", err)
		}

		fmt.Println("✅ MongoDB connection successful!")
		MongoDBClient = client
	})
}

func GetCollection(collectionName string) *mongo.Collection {
	if MongoDBClient == nil {
		log.Fatal("❌ MongoDB is not initialized. Call InitMongoDB() first.")
	}
	return MongoDBClient.Database(MongoDBConfig.Database).Collection(collectionName)
}

func DisconnectMongoDB() {
	if MongoDBClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := MongoDBClient.Disconnect(ctx); err != nil {
			log.Fatalf("❌ Error disconnecting MongoDB: %v", err)
		}
		fmt.Println("🔌 MongoDB connection closed.")
	}
}
