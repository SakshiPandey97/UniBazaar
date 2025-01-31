package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client

func ConnectDB() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default settings")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("MongoDB connection error: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("MongoDB ping error: %v", err)
	}

	log.Println("Connected to MongoDB successfully!")
	DB = client
	return client
}

func GetCollection(collectionName string) *mongo.Collection {
	if DB == nil {
		client := ConnectDB()
		DB = client
	}

	return DB.Database("unibazaar").Collection(collectionName)
}
