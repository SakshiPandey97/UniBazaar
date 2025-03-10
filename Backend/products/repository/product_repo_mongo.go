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

type MongoProductRepository struct {
	collection *mongo.Collection
}

func NewMongoProductRepository() *MongoProductRepository {
	return &MongoProductRepository{
		collection: config.GetCollection("products"),
	}
}

func (repo *MongoProductRepository) getContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func (repo *MongoProductRepository) handleRepoError(err error, context string) error {
	log.Printf("%s: %v\n", context, err)
	return fmt.Errorf("%s: %w", context, err)
}

func (repo *MongoProductRepository) CreateProduct(product model.Product) error {
	log.Printf("Attempting to insert product for UserId: %d\n", product.UserID)

	ctx, cancel := repo.getContextWithTimeout()
	defer cancel()

	_, err := repo.collection.InsertOne(ctx, product)
	if err != nil {
		return repo.handleRepoError(err, "Error inserting product")
	}

	log.Printf("Product inserted successfully for UserId: %d, ProductId: %s\n", product.UserID, product.ProductID)
	return nil
}

func (repo *MongoProductRepository) getProducts(filter bson.M, limit int) ([]model.Product, error) {
	log.Printf("Fetching products with filter: %v, Limit: %d", filter, limit)

	ctx, cancel := repo.getContextWithTimeout()
	defer cancel()

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.D{{Key: "ProductId", Value: 1}}}},
		{{Key: "$limit", Value: int64(limit)}},
	}

	cursor, err := repo.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, repo.handleRepoError(err, "Error fetching products using aggregation")
	}
	defer cursor.Close(ctx)

	var products []model.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, repo.handleRepoError(err, "Error decoding products")
	}

	if len(products) == 0 {
		return nil, repo.handleRepoError(errors.New("no products found"), fmt.Sprintf("Error fetching products"))
	}

	return products, nil
}

func (repo *MongoProductRepository) GetAllProducts(lastID string, limit int) ([]model.Product, error) {
	log.Printf("Fetching products after ID: %s, Limit: %d", lastID, limit)

	filter := bson.M{}
	if lastID != "" {
		log.Printf("Using lastID for filtering: %s", lastID)
		filter = bson.M{"ProductId": bson.M{"$gt": lastID}}
	}

	products, err := repo.getProducts(filter, limit)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (repo *MongoProductRepository) GetProductsByUserID(userID int, lastID string, limit int) ([]model.Product, error) {
	log.Printf("Fetching products for user ID: %d after ID: %s, Limit: %d", userID, lastID, limit)

	filter := bson.M{"UserId": userID}
	if lastID != "" {
		log.Printf("Using lastID for filtering: %s", lastID)
		filter = bson.M{
			"$and": []bson.M{
				{"UserId": userID},
				{"ProductId": bson.M{"$gt": lastID}},
			},
		}
	}

	products, err := repo.getProducts(filter, limit)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (repo *MongoProductRepository) UpdateProduct(userID int, productID string, product model.Product) error {
	log.Printf("Attempting to update product for UserId: %d and ProductId: %s\n", userID, productID)

	ctx, cancel := repo.getContextWithTimeout()
	defer cancel()

	filter := bson.M{"UserId": userID, "ProductId": productID}
	update := bson.M{"$set": product}

	opts := options.Update().SetUpsert(true)

	result, err := repo.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return repo.handleRepoError(err, "Error updating product")
	}

	if result.MatchedCount == 0 && result.UpsertedCount > 0 {
		log.Printf("Product inserted for UserId: %d and ProductId: %s\n", userID, productID)
	} else {
		log.Printf("Product updated successfully for UserId: %d and ProductId: %s\n", userID, productID)
	}

	return nil
}

func (repo *MongoProductRepository) DeleteProduct(userID int, productID string) error {
	log.Printf("Attempting to delete product with ProductID: %s for UserID: %d\n", productID, userID)

	filter := bson.M{
		"UserId":    userID,
		"ProductId": productID,
	}

	result, err := repo.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return repo.handleRepoError(err, "Error deleting product")
	}

	if result.DeletedCount == 0 {
		return errors.New("product not found or already deleted")
	}

	log.Printf("Successfully deleted product with ProductID: %s for UserID: %d\n", productID, userID)
	return nil
}

func (repo *MongoProductRepository) FindProductByUserAndId(userID int, productID string) (*model.Product, error) {
	log.Printf("Attempting to find product for UserId: %d and ProductId: %s\n", userID, productID)

	filter := bson.M{
		"UserId":    userID,
		"ProductId": productID,
	}

	var result model.Product
	ctx, cancel := repo.getContextWithTimeout()
	defer cancel()

	err := repo.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, repo.handleRepoError(errors.New("no products found"), fmt.Sprintf("Product not found for UserId: %d and ProductId: %s", userID, productID))
		}
		return nil, repo.handleRepoError(err, "Error fetching product")
	}

	return &result, nil
}
