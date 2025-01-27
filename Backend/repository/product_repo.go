package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"web-service/config"
	"web-service/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var productCollection *mongo.Collection

func InitProductRepo() {
	productCollection = config.GetCollection("products")
	log.Println("Product repository initialized.")
}

func CreateProduct(product model.Product) error {
	log.Printf("Attempting to create product: %+v\n", product)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := productCollection.InsertOne(ctx, product)
	if err != nil {
		log.Printf("Error creating product: %v\n", err)
		return err
	}

	return nil
}

func GetAllProducts() ([]model.Product, error) {
	log.Println("Fetching all products.")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := productCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error fetching products: %v\n", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	var products []model.Product
	for cursor.Next(ctx) {
		var product model.Product
		if err := cursor.Decode(&product); err != nil {
			log.Printf("Error decoding product: %v\n", err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func GetProductByID(id string) (model.Product, error) {
	log.Printf("Fetching product by ID: %s\n", id)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var product model.Product
	err := productCollection.FindOne(ctx, bson.M{"productId": id}).Decode(&product)

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

func UpdateProduct(id string, updatedProduct model.Product) error {
	log.Printf("Updating product with ID: %s\n", id)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := productCollection.UpdateOne(ctx, bson.M{"productId": id}, bson.M{"$set": updatedProduct})
	if err != nil {
		log.Printf("Error updating product: %v\n", err)
		return err
	}

	if result.MatchedCount == 0 {
		log.Printf("No product found with ID: %s\n", id)
		return errors.New("product not found")
	}

	return nil
}

func DeleteProduct(id string) error {
	log.Printf("Attempting to delete product with ID: %s\n", id)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := productCollection.DeleteOne(ctx, bson.M{"productId": id})
	if err != nil {
		log.Printf("Error deleting product: %v\n", err)
		return err
	}

	if result.DeletedCount == 0 {
		log.Printf("No product found with ID: %s\n", id)
		return errors.New("product not found")
	}

	return nil
}
