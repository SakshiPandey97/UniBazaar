package config

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
)

func TestConnectDB_Success(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	defer os.Unsetenv("MONGO_URI")

	client := ConnectDB()
	if client == nil {
		t.Errorf("Expected a non-nil client, got nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx, nil)
	if err != nil {
		t.Errorf("Ping failed: %v", err)
	}

	err = client.Disconnect(ctx)
	if err != nil {
		log.Printf("Failed to disconnect during test cleanup: %v", err)
	}
}

func TestConnectDB_DefaultURI(t *testing.T) {
	os.Unsetenv("MONGO_URI")

	client := ConnectDB()
	if client == nil {
		t.Errorf("Expected a non-nil client, got nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx, nil)
	if err != nil {
		t.Errorf("Ping failed: %v", err)
	}

	err = client.Disconnect(ctx)
	if err != nil {
		log.Printf("Failed to disconnect during test cleanup: %v", err)
	}
}

func TestGetCollection_Success(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	defer os.Unsetenv("MONGO_URI")

	collectionName := "testCollection"
	collection := GetCollection(collectionName)

	if collection == nil {
		t.Errorf("Expected a non-nil collection, got nil")
	}

	if collection.Name() != collectionName {
		t.Errorf("Expected collection name '%s', got '%s'", collectionName, collection.Name())
	}

	if collection.Database().Name() != "unibazaar" {
		t.Errorf("Expected database name 'unibazaar', got '%s'", collection.Database().Name())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := DB.Disconnect(ctx)
	if err != nil {
		log.Printf("Failed to disconnect during test cleanup: %v", err)
	}
	DB = nil
}

func TestGetCollection_DBNotNil(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	defer os.Unsetenv("MONGO_URI")

	client := ConnectDB()
	DB = client

	collectionName := "testCollection2"
	collection := GetCollection(collectionName)

	if collection == nil {
		t.Errorf("Expected a non-nil collection, got nil")
	}

	if collection.Name() != collectionName {
		t.Errorf("Expected collection name '%s', got '%s'", collectionName, collection.Name())
	}

	if collection.Database().Name() != "unibazaar" {
		t.Errorf("Expected database name 'unibazaar', got '%s'", collection.Database().Name())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := DB.Disconnect(ctx)
	if err != nil {
		log.Printf("Failed to disconnect during test cleanup: %v", err)
	}
	DB = nil
}

func TestConnectDB_ConnectionFailure(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://invalid-uri:27017")
	defer os.Unsetenv("MONGO_URI")

	client := ConnectDB()

	if client == nil {
		t.Errorf("Expected a non-nil client, but got nil")
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := client.Ping(ctx, nil)
		if err == nil {
			t.Errorf("Expected ping to fail, but it succeeded")
		}

		ctxDisconnect, cancelDisconnect := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelDisconnect()

		err = client.Disconnect(ctxDisconnect)
		if err != nil {
			log.Printf("Failed to disconnect during test cleanup: %v", err)
		}

	}
}
