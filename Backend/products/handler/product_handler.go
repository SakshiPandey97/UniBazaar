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

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("Error parsing form data: %v\n", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
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
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
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
	userId := r.URL.Query().Get("userId")
	productId := r.URL.Query().Get("productId")

	if userId == "" || productId == "" {
		http.Error(w, "Missing userId or productId in query parameters", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	log.Printf("Received JSON body: %s", string(body))

	var updateData map[string]interface{}
	if err := json.Unmarshal(body, &updateData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Parsed update data: %+v", updateData)

	if len(updateData) == 0 {
		http.Error(w, "No fields provided for update", http.StatusBadRequest)
		return
	}

	err = repository.UpdateProduct(userId, productId, updateData)
	if err != nil {
		if err.Error() == "product not found" {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error updating product", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product updated successfully"))
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	productId := r.URL.Query().Get("productId")

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
