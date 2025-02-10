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

func CreateProduct(userProduct model.UserProduct) error {
	log.Printf("Attempting to create product: %+v\n", userProduct)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser model.UserProduct
	err := userProductCollection.FindOne(ctx, bson.M{"UserId": userProduct.UserID}).Decode(&existingUser)

	if err == nil {
		update := bson.M{
			"$push": bson.M{"Products": bson.M{
				"ProductId":          userProduct.Products[0].ProductID,
				"ProductTitle":       userProduct.Products[0].ProductTitle,
				"ProductDescription": userProduct.Products[0].ProductDescription,
				"ProductPostDate":    userProduct.Products[0].ProductPostDate,
				"ProductCondition":   userProduct.Products[0].ProductCondition,
				"ProductPrice":       userProduct.Products[0].ProductPrice,
				"ProductLocation":    userProduct.Products[0].ProductLocation,
				"ProductImage":       userProduct.Products[0].ProductImage,
			}},
		}
		_, err = userProductCollection.UpdateOne(ctx, bson.M{"UserId": userProduct.UserID}, update)
		if err != nil {
			log.Printf("Error updating product list: %v\n", err)
			return err
		}

		log.Println("Product added successfully to existing user")
		return nil
	}
	if err == mongo.ErrNoDocuments {
		_, err = userProductCollection.InsertOne(ctx, userProduct)
		if err != nil {
			log.Printf("Error creating new user with product: %v\n", err)
			return err
		}

		log.Println("New user created with product successfully")
		return nil
	}

	log.Printf("Database error: %v\n", err)
	return err
}

func GetAllProducts() ([]model.UserProduct, error) {
	log.Println("Fetching all products.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := userProductCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error fetching products: %v\n", err)
		return nil, err
	}

	defer cursor.Close(ctx)
	var userProducts []model.UserProduct
	if err := cursor.All(ctx, &userProducts); err != nil {
		log.Printf("Error decoding user products: %v\n", err)
		return nil, err
	}

	for i, user := range userProducts {
		for j, product := range user.Products {
			if err == nil {
				preSignedURL, _ := config.GeneratePresignedURL(userProducts[i].Products[j].ProductImage)
				log.Printf("preSignedURL-->%s ", preSignedURL)
				userProducts[i].Products[j].ProductImage = preSignedURL
				log.Printf("preSignedURL 2-->%s ", userProducts[i].Products[j].ProductImage)
			} else {
				log.Printf("Failed to download image for ProductID %s: %v", product.ProductID, err)
			}
		}
	}
	return userProducts, nil
}

func GetProductByID(id string) (model.Product, error) {
	log.Printf("Fetching product by ID: %s\n", id)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product model.Product
	err := userProductCollection.FindOne(ctx, bson.M{"productId": id}).Decode(&product)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Product not found with ID: %s\n", id)
			return product, errors.New("product not found")
		}

		log.Printf("Error fetching product by ID: %v\n", err)
		return product, err
	}

	return product, nil
}

func UpdateProduct(userId, productId string, updateData map[string]interface{}) error {
	userIdStr := fmt.Sprintf("%v", userId)

	updateQuery := bson.M{"$set": bson.M{}}
	for key, value := range updateData {
		updateQuery["$set"].(bson.M)["Products.$."+key] = value
	}

	filter := bson.M{
		"UserId":             userIdStr,
		"Products.ProductId": productId,
	}

	log.Printf("Filter: %+v", filter)
	log.Printf("Update Query: %+v", updateQuery)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := userProductCollection.UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		log.Printf("Error updating product: %v\n", err)
		return errors.New("failed to update product")
	}

	if result.MatchedCount == 0 {
		log.Printf("No matching product found for UserId: %s and ProductId: %s\n", userId, productId)
		return errors.New("product not found")
	}

	log.Printf("Product with ID %s updated successfully\n", productId)
	return nil
}

func DeleteProduct(userId, productId string) error {
	log.Printf("Attempting to delete product with ProductID: %s for UserID: %s\n", productId, userId)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"UserId": userId,
		"Products": bson.M{
			"$elemMatch": bson.M{
				"ProductId": productId,
			},
		},
	}

	result, err := userProductCollection.UpdateOne(
		ctx,
		filter,
		bson.M{
			"$pull": bson.M{"Products": bson.M{"ProductId": productId}},
		},
	)

	if err != nil {
		log.Printf("Error deleting product: %v\n", err)
		return fmt.Errorf("failed to delete product with ProductID %s for UserID %s: %w", productId, userId, err)
	}

	if result.ModifiedCount == 0 {
		log.Printf("No product found with ProductID: %s for UserID: %s\n", productId, userId)
		return errors.New("product not found or already deleted")
	}

	log.Printf("Successfully deleted product with ProductID: %s for UserID: %s\n", productId, userId)
	return nil
}
