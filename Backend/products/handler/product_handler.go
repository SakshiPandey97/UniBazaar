package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"web-service/helper"
	"web-service/model"
	"web-service/repository"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create a new product.")
	log.Printf(r.FormValue("productImage"))
	// Parse the form data (max 10MB file size)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing form data: %v\n", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Extract text fields
	product := model.Product{
		ProductID:          uuid.NewString(),
		ProductTitle:       r.FormValue("productTitle"),
		ProductDescription: r.FormValue("productDescription"),
		ProductPostDate:    r.FormValue("productPostDate"),
		ProductCondition:   0,
		ProductPrice:       0.0,
		ProductLocation:    r.FormValue("productLocation"),
	}

	// Convert form values to appropriate types
	fmt.Sscanf(r.FormValue("productCondition"), "%d", &product.ProductCondition)
	fmt.Sscanf(r.FormValue("productPrice"), "%f", &product.ProductPrice)

	// Extract the file
	file, _, err := r.FormFile("productImage")
	if err != nil {
		log.Printf("Error retrieving file: %v\n", err)
		http.Error(w, "File upload error", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		log.Printf("Error reading file: %v\n", err)
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}
	// Store the file reference in the product object
	s3IamgeKey, uploadError := helper.UploadToS3Bucket(product.ProductID, r.FormValue("userId"), buf.Bytes(), r.FormValue("productImageType"))
	if uploadError != nil {
		log.Printf("Error uploading file to S3: %v\n", err)
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}
	product.ProductImage = s3IamgeKey
	userProduct := model.UserProduct{
		UserID:   r.FormValue("userId"),
		Products: []model.Product{product},
	}
	if err := repository.CreateProduct(userProduct); err != nil {
		log.Printf("Error creating product: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	log.Printf("Product created successfully: %+v\n", product)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to fetch all products.")

	products, err := repository.GetAllProducts()
	if err != nil {
		log.Printf("Error fetching products: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d products.\n", len(products))
	jsonData, _ := json.MarshalIndent(products, "", "  ")
	log.Printf("Final JSON Output:\n%s", string(jsonData))
	json.NewEncoder(w).Encode(products)
}

func GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Received request to fetch product with ID: %s\n", id)

	product, err := repository.GetProductByID(id)
	if err != nil {
		log.Printf("Error fetching product with ID %s: %v\n", id, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Product found: %+v\n", product)
	json.NewEncoder(w).Encode(product)
}

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	log.Printf("Received request to update product with ID: %s\n", id)

	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := product.Validate(); err != nil {
		log.Printf("Validation error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := product.Validate(); err != nil {
		log.Printf("Validation error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Product with ID %s updated successfully.\n", id)
	w.WriteHeader(http.StatusOK)
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	productId := r.URL.Query().Get("productId")

	// Validate query parameters
	if userId == "" || productId == "" {
		http.Error(w, "Missing userId or productId in query parameters", http.StatusBadRequest)
		return
	}

	log.Printf("Received request to delete product with ID: %s by user %s\n", productId, userId)

	if err := repository.DeleteProduct(userId, productId); err != nil {
		log.Printf("Error deleting product with ID %s: %v\n", productId, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Product with ID %s deleted successfully.\n", productId)
	w.WriteHeader(http.StatusNoContent)
}
