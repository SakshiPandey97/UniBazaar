package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"web-service/config"
	"web-service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userProductCollection *mongo.Collection

func InitProductRepo() {
	userProductCollection = config.GetCollection("products")
	log.Println("Product repository initialized.")
}

func getContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func handleRepoError(err error, context string) error {
	log.Printf("%s: %v\n", context, err)
	return fmt.Errorf("%s: %w", context, err)
}

func CreateProduct(product model.Product) error {
	log.Printf("Attempting to insert product for UserId: %d\n", product.UserID)

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	_, err := userProductCollection.InsertOne(ctx, product)
	if err != nil {
		return handleRepoError(err, "Error inserting product")
	}

	log.Printf("Product inserted successfully for UserId: %d, ProductId: %s\n", product.UserID, product.ProductID)
	return nil
}

func GetAllProducts() ([]model.Product, error) {
	log.Println("Fetching all products.")

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	cursor, err := userProductCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, handleRepoError(err, "Error fetching products")
	}
	defer cursor.Close(ctx)

	var products []model.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, handleRepoError(err, "Error decoding products")
	}

	return products, nil
}

func GetProductsByUserID(userID int) ([]model.Product, error) {
	log.Printf("Fetching products for user ID: %d\n", userID)

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	cursor, err := userProductCollection.Find(ctx, bson.M{"UserId": userID})
	if err != nil {
		return nil, handleRepoError(err, "Error fetching products for user")
	}
	defer cursor.Close(ctx)

	var products []model.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, handleRepoError(err, "Error decoding user products")
	}

	if len(products) == 0 {
		return nil, fmt.Errorf("no products found for user ID: %d", userID)
	}

	return products, nil
}

func UpdateProduct(userId int, productId string, product model.Product) error {
	log.Printf("Attempting to update product for UserId: %d and ProductId: %s\n", userId, productId)

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	filter := bson.M{"UserId": userId, "ProductId": productId}
	update := bson.M{"$set": product}

	opts := options.Update().SetUpsert(true)

	result, err := userProductCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return handleRepoError(err, "Error updating product")
	}

	if result.MatchedCount == 0 && result.UpsertedCount > 0 {
		log.Printf("Product inserted for UserId: %d and ProductId: %s\n", userId, productId)
	} else {
		log.Printf("Product updated successfully for UserId: %d and ProductId: %s\n", userId, productId)
	}

	return nil
}

func DeleteProduct(userId int, productId string) error {
	log.Printf("Attempting to delete product with ProductID: %s for UserID: %d\n", productId, userId)

	filter := bson.M{
		"UserId":    userId,
		"ProductId": productId,
	}

	result, err := userProductCollection.DeleteOne(
		context.Background(),
		filter,
	)

	if err != nil {
		return handleRepoError(err, "Error deleting product")
	}

	if result.DeletedCount == 0 {
		return errors.New("product not found or already deleted")
	}

	log.Printf("Successfully deleted product with ProductID: %s for UserID: %d\n", productId, userId)
	return nil
}

func FindProductByUserAndId(userId int, productId string) (*model.Product, error) {
	log.Printf("Attempting to find product for UserId: %d and ProductId: %s\n", userId, productId)

	filter := bson.M{
		"UserId":    userId,
		"ProductId": productId,
	}

	var result model.Product
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	err := userProductCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("product not found for UserId: %d and ProductId: %s", userId, productId)
		}
		return nil, handleRepoError(err, "Error fetching product")
	}

	return &result, nil
}
