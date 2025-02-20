package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"web-service/helper"
	"web-service/model"
	"web-service/repository"

	"github.com/gorilla/mux"
)

// @Summary Create a new product
// @Description Creates a new product by parsing form data, uploading images to S3, and saving it to the database. The product is linked to the user via their User ID.
// @Tags Products
// @Accept multipart/form-data
// @Produce json
// @Param UserId formData int true "User ID (form data)"
// @Param productTitle formData string true "Product title"
// @Param productDescription formData string false "Product description"
// @Param productPrice formData float64 true "Product price"
// @Param productCondition formData int true "Product condition"
// @Param productLocation formData string true "Product location"
// @Param productImage formData file true "Product image"
// @Success 201 {object} model.Product "Product created successfully"
// @Failure 400 {object} model.ErrorResponse "Invalid User ID or form data"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /products [post]
func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create a new product.")

	userID, err := helper.GetUserID(r.FormValue("userId"))
	if err != nil {
		handleError(w, "Invalid userId", err, http.StatusBadRequest)
		return
	}

	product, err := helper.ParseFormAndCreateProduct(r, userID)
	if err != nil {
		handleError(w, "Error creating product", err, http.StatusBadRequest)
		return
	}

	product.UserID = userID

	s3ImageKey, err := handleProductImageUpload(w, r, &product)
	if err != nil {
		return
	}
	product.ProductImage = s3ImageKey

	if err := repository.CreateProduct(product); err != nil {
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

// @Summary Get all products in the system
// @Description Fetch all products from the system, regardless of the user ID. If no products are found, an error is returned.
// @Tags Products
// @Accept json
// @Produce json
// @Success 200 {array} model.Product "List of all products"
// @Failure 404 {object} model.ErrorResponse "No products found in the system"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /products [get]
func GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to fetch all products.")

	products, err := repository.GetAllProducts()
	if err != nil {
		handleError(w, "Error fetching products", err, http.StatusInternalServerError)
		return
	}

	if len(products) == 0 {
		handleError(w, "No products found in the system", fmt.Errorf("no products found"), http.StatusNotFound)
		return
	}

	products = repository.GetPreSignedURLs(products)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		handleError(w, "Error encoding response", err, http.StatusInternalServerError)
		return
	}
}

// @Summary Get all products for a specific user by user ID
// @Description Fetch all products listed by a user, identified by their user ID. If no products are found, an error is returned.
// @Tags Products
// @Accept json
// @Produce json
// @Param UserId path int true "User ID"
// @Success 200 {array} model.Product "List of products"
// @Failure 400 {object} model.ErrorResponse "Invalid user ID"
// @Failure 404 {object} model.ErrorResponse "No products found for the given user ID"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /products/{UserId} [get]
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

// @Summary Update a product by user ID and product ID
// @Description Update a product's details based on the user ID and product ID. The product image is also updated if provided.
// @Tags Products
// @Accept json
// @Produce json
// @Param UserId path int true "User ID"
// @Param ProductId path string true "Product ID"
// @Param product body model.Product true "Product Details"
// @Success 200 {object} model.Product "Updated product"
// @Failure 400 {object} model.ErrorResponse "Invalid request"
// @Failure 404 {object} model.ErrorResponse "Product not found"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /products/{UserId}/{ProductId} [put]
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

	updatedProduct, err := helper.ParseFormAndCreateProduct(r, userId)
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

// @Summary Delete a product by user ID and product ID
// @Description Delete a product from the system based on the user ID and product ID. This also removes the associated image from S3 if available.
// @Tags Products
// @Param UserId path int true "User ID"
// @Param ProductId path string true "Product ID"
// @Success 204 "Product deleted"
// @Failure 400 {object} model.ErrorResponse "Invalid request"
// @Failure 404 {object} model.ErrorResponse "Product not found"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /products/{UserId}/{ProductId} [delete]
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
