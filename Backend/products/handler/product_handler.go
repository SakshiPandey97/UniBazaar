package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"web-service/helper"
	"web-service/model"
	"web-service/repository"

	"github.com/gorilla/mux"
)

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create a new product.")

	UserID, err := helper.GetUserID(r.FormValue("userId"))
	if err != nil {
		handleError(w, "Invalid userId", err, http.StatusBadRequest)
		return
	}

	product, err := helper.ParseFormAndCreateProduct(r)
	if err != nil {
		handleError(w, "Error creating product", err, http.StatusBadRequest)
		return
	}

	s3ImageKey, err := handleProductImageUpload(w, r, &product)
	if err != nil {
		return
	}
	product.ProductImage = s3ImageKey

	userProduct := model.UserProduct{
		UserID:   UserID,
		Products: []model.Product{product},
	}

	if err := repository.CreateProduct(userProduct); err != nil {
		handleError(w, "Error creating product in database", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(product); err != nil {
		handleError(w, "Error encoding response", err, http.StatusInternalServerError)
		return
	}
}

func GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to fetch all products.")

	userProducts, err := repository.GetAllProducts()
	if err != nil {
		handleError(w, "Error fetching products", err, http.StatusInternalServerError)
		return
	}

	for i := range userProducts {
		userProducts[i].Products = repository.GetPreSignedURLs(userProducts[i].Products)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userProducts); err != nil {
		handleError(w, "Error encoding response", err, http.StatusInternalServerError)
		return
	}
}

func GetAllProductsByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := helper.GetUserID(mux.Vars(r)["UserId"])
	if err != nil {
		handleError(w, "Invalid userId", err, http.StatusBadRequest)
		return
	}

	log.Printf("Received request to fetch all products for user ID: %d\n", userID)

	products, err := repository.GetProductsByUserID(userID)
	if err != nil {
		handleError(w, "Error fetching products for user", err, http.StatusNotFound)
		return
	}

	products = repository.GetPreSignedURLs(products)

	log.Printf("Found %d products for user ID %d\n", len(products), userID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		handleError(w, "Error encoding response", err, http.StatusInternalServerError)
		return
	}
}

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetUserID(mux.Vars(r)["UserId"])
	if err != nil {
		handleError(w, "Invalid or missing userId", err, http.StatusBadRequest)
		return
	}

	productId := mux.Vars(r)["ProductId"]
	if productId == "" {
		handleError(w, "Missing productId in URL parameters", nil, http.StatusBadRequest)
		return
	}

	existingProduct, err := repository.FindProductByUserAndId(userId, productId)
	if err != nil {
		handleError(w, "Error fetching product", err, http.StatusNotFound)
		return
	}

	updatedProduct, err := helper.ParseFormAndCreateProduct(r)
	if err != nil {
		handleError(w, "Error parsing form data", err, http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var imageDeleteErr, imageUploadErr error
	var newS3ImageKey string

	go func() {
		defer wg.Done()
		if existingProduct.ProductImage != "" {
			imageDeleteErr = repository.DeleteImageFromS3(existingProduct.ProductImage)
		}
	}()

	go func() {
		defer wg.Done()
		newS3ImageKey, imageUploadErr = handleProductImageUpload(w, r, &updatedProduct)
	}()

	wg.Wait()

	if imageDeleteErr != nil {
		log.Printf("Error deleting old image: %v", imageDeleteErr)
	}

	if imageUploadErr != nil {
		handleError(w, "Error uploading new image", imageUploadErr, http.StatusInternalServerError)
		return
	}

	updatedProduct.ProductImage = newS3ImageKey

	err = repository.UpdateProduct(userId, productId, updatedProduct)
	if err != nil {
		handleError(w, "Error updating product", err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedProduct); err != nil {
		handleError(w, "Error encoding response", err, http.StatusInternalServerError)
	}
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetUserID(mux.Vars(r)["UserId"])
	if err != nil {
		handleError(w, "Invalid or missing userId", err, http.StatusBadRequest)
		return
	}

	productId := mux.Vars(r)["ProductId"]
	if productId == "" {
		handleError(w, "Missing productId in path parameters", errors.New("productId is required"), http.StatusBadRequest)
		return
	}

	log.Printf("Received request to delete product with ID: %s by user %d\n", productId, userId)

	product, err := repository.FindProductByUserAndId(userId, productId)
	if err != nil {
		handleError(w, "Error fetching product", err, http.StatusNotFound)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var imageDeleteErr, dbDeleteErr error

	go func() {
		defer wg.Done()
		if product.ProductImage != "" {
			imageDeleteErr = repository.DeleteImageFromS3(product.ProductImage)
		}
	}()

	go func() {
		defer wg.Done()
		dbDeleteErr = repository.DeleteProduct(userId, productId)
	}()

	wg.Wait()

	if imageDeleteErr != nil {
		log.Printf("Error deleting image from S3: %v\n", imageDeleteErr)
	}

	if dbDeleteErr != nil {
		handleError(w, "Error deleting product", dbDeleteErr, http.StatusInternalServerError)
		return
	}

	log.Printf("Product with ID %s deleted successfully.\n", productId)
	w.WriteHeader(http.StatusNoContent)
}

func handleProductImageUpload(w http.ResponseWriter, r *http.Request, product *model.Product) (string, error) {
	imageData, format, err := helper.ParseProductImage(r)
	if err != nil {
		handleError(w, "Error reading image", err, http.StatusBadRequest)
		return "", err
	}

	s3ImageKey, err := repository.UploadToS3Bucket(product.ProductID, r.FormValue("userId"), imageData.Bytes(), format)
	if err != nil {
		handleError(w, "Error uploading image to S3", err, http.StatusInternalServerError)
		return "", err
	}
	return s3ImageKey, nil
}

func handleError(w http.ResponseWriter, message string, err error, statusCode int) {
	log.Printf("%s: %v\n", message, err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := model.ErrorResponse{
		Error:   message,
		Details: err.Error(),
	}

	json.NewEncoder(w).Encode(response)
}
