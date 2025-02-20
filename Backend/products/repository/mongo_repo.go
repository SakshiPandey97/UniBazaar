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
)

var userProductCollection *mongo.Collection

func InitProductRepo() {
	userProductCollection = config.GetCollection("user_product")
	log.Println("Product repository initialized.")
}

func getContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func handleRepoError(err error, context string) error {
	log.Printf("%s: %v\n", context, err)
	return fmt.Errorf("%s: %w", context, err)
}

func findUserByID(userID int) (model.UserProduct, error) {
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	var userProduct model.UserProduct
	err := userProductCollection.FindOne(ctx, bson.M{"UserId": userID}).Decode(&userProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.UserProduct{}, errors.New("user not found")
		}
		return model.UserProduct{}, err
	}
	return userProduct, nil
}

func CreateProduct(userProduct model.UserProduct) error {
	log.Printf("Attempting to create product for UserId: %d\n", userProduct.UserID)

	_, err := findUserByID(userProduct.UserID)
	if err == nil {
		update := bson.M{
			"$push": bson.M{"Products": userProduct.Products[0]},
		}
		ctx, cancel := getContextWithTimeout()
		defer cancel()

		_, err := userProductCollection.UpdateOne(ctx, bson.M{"UserId": userProduct.UserID}, update)
		if err != nil {
			return handleRepoError(err, "Error updating product list")
		}

		log.Printf("Product added successfully to existing user with UserId: %d\n", userProduct.UserID)
		return nil
	}

	if err == mongo.ErrNoDocuments {
		ctx, cancel := getContextWithTimeout()
		defer cancel()

		_, err = userProductCollection.InsertOne(ctx, userProduct)
		if err != nil {
			return handleRepoError(err, "Error creating new user with product")
		}

		log.Printf("New user created with product successfully for UserId: %d\n", userProduct.UserID)
		return nil
	}

	return handleRepoError(err, "Database error while handling UserId")
}

func GetAllProducts() ([]model.UserProduct, error) {
	log.Println("Fetching all products.")

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	cursor, err := userProductCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, handleRepoError(err, "Error fetching products")
	}
	defer cursor.Close(ctx)

	var userProducts []model.UserProduct
	if err := cursor.All(ctx, &userProducts); err != nil {
		return nil, handleRepoError(err, "Error decoding user products")
	}

	return userProducts, nil
}

func GetProductsByUserID(userID int) ([]model.Product, error) {
	log.Printf("Fetching products for user ID: %d\n", userID)

	userProduct, err := findUserByID(userID)
	if err != nil {
		return nil, handleRepoError(err, "Error fetching products for user")
	}

	return userProduct.Products, nil
}

func UpdateProduct(userId int, productId string, product model.Product) error {
	log.Printf("Attempting to update product for UserId: %d\n", userId)

	updateData := bson.M{
		"$set": bson.M{
			"Products.$": product,
		},
	}

	filter := bson.M{
		"UserId":             userId,
		"Products.ProductId": productId,
	}

	ctx, cancel := getContextWithTimeout()
	defer cancel()

	result, err := userProductCollection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		return handleRepoError(err, "Error updating product")
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("product not found for UserId: %d and ProductId: %s", userId, productId)
	}

	log.Printf("Product with ID %s updated successfully\n", productId)
	return nil
}

func DeleteProduct(userId int, productId string) error {
	log.Printf("Attempting to delete product with ProductID: %s for UserID: %d\n", productId, userId)

	filter := bson.M{
		"UserId": userId,
		"Products": bson.M{
			"$elemMatch": bson.M{
				"ProductId": productId,
			},
		},
	}

	result, err := userProductCollection.UpdateOne(
		context.Background(),
		filter,
		bson.M{
			"$pull": bson.M{"Products": bson.M{"ProductId": productId}},
		},
	)

	if err != nil {
		return handleRepoError(err, "Error deleting product")
	}

	if result.ModifiedCount == 0 {
		return errors.New("product not found or already deleted")
	}

	log.Printf("Successfully deleted product from MongoDB with ProductID: %s for UserID: %d\n", productId, userId)
	return nil
}

func FindProductByUserAndId(userId int, productId string) (*model.Product, string, error) {
	log.Printf("Attempting to find product for UserId: %d and ProductId: %s\n", userId, productId)

	filter := bson.M{
		"UserId":             userId,
		"Products.ProductId": productId,
	}

	var userProduct model.UserProduct
	ctx, cancel := getContextWithTimeout()
	defer cancel()

	err := userProductCollection.FindOne(ctx, filter).Decode(&userProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, "", fmt.Errorf("product not found for UserId: %d and ProductId: %s", userId, productId)
		}
		return nil, "", handleRepoError(err, "Error fetching product")
	}

	for _, product := range userProduct.Products {
		if product.ProductID == productId {
			return &product, product.ProductImage, nil
		}
	}

	return nil, "", fmt.Errorf("product not found in user's list")
}
