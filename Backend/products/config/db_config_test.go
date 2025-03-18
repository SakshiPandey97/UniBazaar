package config

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConnectDB_Success(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	defer os.Unsetenv("MONGO_URI")

	client, err := ConnectDB()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	assert.NoError(t, err)

	err = client.Disconnect(ctx)
	assert.NoError(t, err)
}

func TestConnectDB_DefaultURI(t *testing.T) {
	os.Unsetenv("MONGO_URI")

	client, err := ConnectDB()
	assert.NoError(t, err)
	assert.NotNil(t, client)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	assert.NoError(t, err)

	err = client.Disconnect(ctx)
	assert.NoError(t, err)
}

func TestGetCollection_Success(t *testing.T) {
	t.Run("isolated test", func(t *testing.T) {
		os.Setenv("MONGO_URI", "mongodb://localhost:27017")
		defer os.Unsetenv("MONGO_URI")

		DB = nil

		collectionName := "testCollection"
		collection, err := GetCollection(collectionName)
		assert.NoError(t, err)
		assert.NotNil(t, collection)

		assert.Equal(t, collectionName, collection.Name())
		assert.Equal(t, "unibazaar", collection.Database().Name())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if DB != nil {
			err = DB.Disconnect(ctx)
			assert.NoError(t, err)
			DB = nil
		}
	})
}

func TestGetCollection_DBNotNil(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://localhost:27017")
	defer os.Unsetenv("MONGO_URI")

	client, err := ConnectDB()
	assert.NoError(t, err)
	DB = client

	collectionName := "testCollection2"
	collection, err := GetCollection(collectionName)
	assert.NoError(t, err)
	assert.NotNil(t, collection)

	assert.Equal(t, collectionName, collection.Name())
	assert.Equal(t, "unibazaar", collection.Database().Name())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = DB.Disconnect(ctx)
	assert.NoError(t, err)
	DB = nil
}

func TestConnectDB_ConnectionFailure(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://invalid-uri:27017")
	defer os.Unsetenv("MONGO_URI")

	client, err := ConnectDB()
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestGetCollection_ConnectDBFailure(t *testing.T) {
	os.Setenv("MONGO_URI", "mongodb://invalid-uri:27017")
	defer os.Unsetenv("MONGO_URI")

	collectionName := "testCollection"
	collection, err := GetCollection(collectionName)

	assert.Error(t, err)
	assert.Nil(t, collection)
}
