package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"web-service/config"
	"web-service/helper"
	"web-service/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

var userProductCollection *mongo.Collection

func InitProductRepo() {
	userProductCollection = config.GetCollection("user_product")
	log.Println("Product repository initialized.")
}

func CreateProduct(userProduct model.UserProduct, r *http.Request) error {
	log.Printf("Attempting to create product: %+v\n", userProduct)

	product, err := BuildProductFromRequest(r)
	if err != nil {
		return err
	}

	s3ImageKey, err := HandleProductImage(r, product.ProductID, r.FormValue("userId"))
	if err != nil {
		return err
	}

	product.ProductImage = s3ImageKey

	userProduct.UserID = r.FormValue("userId")
	userProduct.Products = append(userProduct.Products, product)

	err = CreateOrUpdateProduct(userProduct, userProductCollection)
	if err != nil {
		log.Printf("Error handling product creation/update: %v\n", err)
		return err
	}

	log.Printf("Product created successfully: %+v\n", product)
	return nil
}

func BuildProductFromRequest(r *http.Request) (model.Product, error) {
	log.Println("Building product from request.")

	err := r.ParseMultipartForm(10 << 20) // 10MB max size
	if err != nil {
		log.Printf("Error parsing form data: %v\n", err)
		return model.Product{}, errors.New("failed to parse form data")
	}

	product := model.Product{
		ProductID:          uuid.NewString(),
		ProductTitle:       r.FormValue("productTitle"),
		ProductDescription: r.FormValue("productDescription"),
		ProductPostDate:    r.FormValue("productPostDate"),
		ProductCondition:   0,
		ProductPrice:       0.0,
		ProductLocation:    r.FormValue("productLocation"),
	}

	fmt.Sscanf(r.FormValue("productCondition"), "%d", &product.ProductCondition)
	fmt.Sscanf(r.FormValue("productPrice"), "%f", &product.ProductPrice)

	return product, nil
}

func HandleProductImage(r *http.Request, productID string, userID string) (string, error) {
	log.Println("Handling product image.")

	file, _, err := r.FormFile("productImage")
	if err != nil {
		log.Printf("Error retrieving file: %v\n", err)
		return "", errors.New("file upload error")
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		log.Printf("Error reading file: %v\n", err)
		return "", errors.New("error processing file")
	}

	s3ImageKey, uploadError := helper.UploadToS3Bucket(productID, userID, buf.Bytes(), r.FormValue("productImageType"))
	if uploadError != nil {
		log.Printf("Error uploading file to S3: %v\n", uploadError)
		return "", errors.New("error uploading file to S3")
	}

	return s3ImageKey, nil
}

func GetAllProducts() ([]model.UserProduct, error) {
	log.Println("Fetching all products.")

	userProducts, err := GetAllProductsFromDB(userProductCollection)
	if err != nil {
		log.Printf("Error fetching products from DB: %v\n", err)
		return nil, err
	}

	for i, user := range userProducts {
		for j, product := range user.Products {
			if err == nil {
				preSignedURL, _ := config.GeneratePresignedURL(userProducts[i].Products[j].ProductImage)
				userProducts[i].Products[j].ProductImage = preSignedURL
			} else {
				log.Printf("Failed to download image for ProductID %s: %v", product.ProductID, err)
			}
		}
	}

	return userProducts, nil
}

func GetProductByID(id string) (model.Product, error) {
	log.Printf("Fetching product by ID: %s\n", id)

	product, err := GetProductByIDFromDB(id, userProductCollection)
	if err != nil {
		log.Printf("Error fetching product by ID from DB: %v\n", err)
		return product, err
	}

	return product, nil
}

func UpdateProduct(w http.ResponseWriter, r *http.Request, userId, productId string) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return err
	}
	log.Printf("Received JSON body: %s", string(body))

	var updateData map[string]interface{}
	if err := json.Unmarshal(body, &updateData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return err
	}

	log.Printf("Parsed update data: %+v", updateData)

	if len(updateData) == 0 {
		http.Error(w, "No fields provided for update", http.StatusBadRequest)
		return nil
	}

	err = UpdateProductInDB(userId, productId, updateData, userProductCollection)
	if err != nil {
		if err.Error() == "product not found" {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error updating product", http.StatusInternalServerError)
		}
		return err
	}

	return nil
}

func DeleteProduct(userId, productId string) error {
	log.Printf("Attempting to delete product with ProductID: %s for UserID: %s\n", productId, userId)

	err := DeleteProductInDB(userId, productId, userProductCollection)
	if err != nil {
		log.Printf("Error deleting product: %v\n", err)
		return err
	}

	return nil
}
