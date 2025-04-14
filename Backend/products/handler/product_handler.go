package handler

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"web-service/helper"
	"web-service/model"
	"web-service/repository"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	ProductRepo repository.ProductRepository
	ImageRepo   repository.ImageRepository
}

func NewProductHandler(productRepo repository.ProductRepository, imageRepo repository.ImageRepository) *ProductHandler {
	return &ProductHandler{
		ProductRepo: productRepo,
		ImageRepo:   imageRepo,
	}
}

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
func (h *ProductHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create a new product.")

	userID, err := helper.GetUserID(r.FormValue("userId"))
	if err != nil {
		HandleError(w, err, "Invalid userId")
		return
	}

	product, err := helper.ParseFormAndCreateProduct(r, userID)
	if err != nil {
		HandleError(w, err, "Error creating product")
		return
	}
	product.UserID = userID

	s3ImageKey, err := h.handleProductImageUpload(w, r, &product)
	if err != nil {
		return
	}

	product.ProductImage = s3ImageKey

	if err := h.ProductRepo.CreateProduct(product); err != nil {
		HandleError(w, err, "Error creating product")
		return
	}

	HandleSuccessResponse(w, http.StatusCreated, product)
}

// @Summary Get all products in the system
// @Description Fetch all products from the system, regardless of the user ID. If no products are found, an error is returned.
// @Tags Products
// @Accept json
// @Produce json
// @Param lastId query string false "ID of the last product to fetch" required=false
// @Param limit query int false "Number of products to fetch (default is 10)" required=false
// @Success 200 {array} model.Product "List of all products"
// @Failure 404 {object} model.ErrorResponse "No products found in the system"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /products [get]
func (h *ProductHandler) GetAllProductsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to fetch all products.")

	lastID := r.URL.Query().Get("lastId")
	limitStr := r.URL.Query().Get("limit")

	limit := helper.ParseLimit(limitStr)

	products, err := h.ProductRepo.GetAllProducts(lastID, limit)
	if err != nil {
		HandleError(w, err, "Error fetching products")
		return
	}

	products = h.ImageRepo.GetPreSignedURLs(products)

	HandleSuccessResponse(w, http.StatusOK, products)
}

// @Summary Get all products for a specific user by user ID
// @Description Fetch all products listed by a user, identified by their user ID. If no products are found, an error is returned.
// @Tags Products
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param lastID query string false "ID of the last product to fetch" required=false
// @Param limit query int false "Number of products to fetch (default is 10)" required=false
// @Success 200 {array} model.Product "List of products"
// @Failure 400 {object} model.ErrorResponse "Invalid user ID"
// @Failure 404 {object} model.ErrorResponse "No products found for the given user ID"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /products/{userId} [get]
func (h *ProductHandler) GetAllProductsByUserIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := helper.GetUserID(mux.Vars(r)["UserId"])
	if err != nil {
		HandleError(w, err, "Invalid userId")
		return
	}

	lastID := r.URL.Query().Get("lastId")
	limitStr := r.URL.Query().Get("limit")

	limit := helper.ParseLimit(limitStr)

	log.Printf("Received request to fetch all products for user ID: %d with lastID: %s and limit: %d\n", userID, lastID, limit)

	products, err := h.ProductRepo.GetProductsByUserID(userID, lastID, limit)
	if err != nil {
		HandleError(w, err, "Error fetching products")
		return
	}

	products = h.ImageRepo.GetPreSignedURLs(products)

	log.Printf("Found %d products for user ID %d\n", len(products), userID)

	HandleSuccessResponse(w, http.StatusOK, products)
}

// @Summary Update a product by user ID and product ID
// @Description Update a product's details based on the user ID and product ID. The product image is also updated if provided.
// @Tags Products
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param productId path string true "Product ID"
// @Param product body model.Product true "Product Details"
// @Success 200 {object} model.Product "Updated product"
// @Failure 400 {object} model.ErrorResponse "Invalid request"
// @Failure 404 {object} model.ErrorResponse "Product not found"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /products/{userId}/{productId} [put]
func (h *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetUserID(mux.Vars(r)["UserId"])
	if err != nil {
		HandleError(w, err, "Invalid userId")
		return
	}

	productId, err := helper.CheckParam(mux.Vars(r)["ProductId"])
	if err != nil {
		HandleError(w, err, "Error checking product ID")
		return
	}

	existingProduct, err := h.ProductRepo.FindProductByUserAndId(userId, productId)
	if err != nil {
		HandleError(w, err, "Error finding existing product")
		return
	}

	updatedProduct, err := helper.ParseFormAndCreateProduct(r, userId)
	if err != nil {
		HandleError(w, err, "Error parsing form data")
		return
	}

	_, _, err = r.FormFile("productImage")
	if err == http.ErrMissingFile {
		updatedProduct.ProductImage = existingProduct.ProductImage
	} else if err != nil {
		HandleError(w, err, "Error retrieving uploaded image")
		return
	} else {
		var wg sync.WaitGroup
		wg.Add(2)
		var imageDeleteErr, imageUploadErr error
		var newS3ImageKey string

		go func() {
			defer wg.Done()
			if existingProduct.ProductImage != "" {
				imageDeleteErr = h.ImageRepo.DeleteImage(existingProduct.ProductImage)
			}
		}()

		go func() {
			defer wg.Done()
			newS3ImageKey, imageUploadErr = h.handleProductImageUpload(w, r, &updatedProduct)
		}()

		wg.Wait()

		if imageDeleteErr != nil {
			log.Printf("Error deleting old image: %v", imageDeleteErr)
		}

		if imageUploadErr != nil {
			HandleError(w, imageUploadErr, "Error uploading image")
			return
		}

		updatedProduct.ProductImage = newS3ImageKey
	}

	err = h.ProductRepo.UpdateProduct(userId, productId, updatedProduct)
	if err != nil {
		HandleError(w, err, "Error updating product in database")
		return
	}

	productsWithURL := h.ImageRepo.GetPreSignedURLs([]model.Product{updatedProduct})
	if len(productsWithURL) > 0 {
		HandleSuccessResponse(w, http.StatusOK, productsWithURL[0])
	} else {
		HandleError(w, errors.New("failed to generate pre-signed URL"), "Image URL generation failed")
	}
}

// @Summary Delete a product by user ID and product ID
// @Description Delete a product from the system based on the user ID and product ID. This also removes the associated image from S3 if available.
// @Tags Products
// @Param userId path int true "User ID"
// @Param productId path string true "Product ID"
// @Success 204 "Product deleted"
// @Failure 400 {object} model.ErrorResponse "Invalid request"
// @Failure 404 {object} model.ErrorResponse "Product not found"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /products/{userId}/{productId} [delete]
func (h *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := helper.GetUserID(mux.Vars(r)["UserId"])
	if err != nil {
		HandleError(w, err, "Invalid userId")
		return
	}

	productId, err := helper.CheckParam(mux.Vars(r)["ProductId"])
	if err != nil {
		HandleError(w, err, "Error checking product ID")
		return
	}

	log.Printf("Received request to delete product with ID: %s by user %d\n", productId, userId)

	product, err := h.ProductRepo.FindProductByUserAndId(userId, productId)
	if err != nil {
		HandleError(w, err, "Error updating product")
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	var imageDeleteErr, dbDeleteErr error

	go func() {
		defer wg.Done()
		if product.ProductImage != "" {
			imageDeleteErr = h.ImageRepo.DeleteImage(product.ProductImage)
		}
	}()

	go func() {
		defer wg.Done()
		dbDeleteErr = h.ProductRepo.DeleteProduct(userId, productId)
	}()

	wg.Wait()

	if imageDeleteErr != nil {
		HandleError(w, imageDeleteErr, "Error deleting image")
		return
	}

	if dbDeleteErr != nil {
		HandleError(w, dbDeleteErr, "Error deleting product")
		return
	}

	log.Printf("Product with ID %s deleted successfully.\n", productId)
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Search products
// @Description Searches products based on a query and optional limit.
// @Tags Products
// @Param query query string true "Search query"
// @Param limit query int false "Limit the number of results"
// @Success 200 {object} []model.Product "List of products matching the search query"
// @Failure 400 {object} model.ErrorResponse "Invalid request or missing query parameter"
// @Failure 404 {object} model.ErrorResponse "No products found for the given query"
// @Failure 500 {object} model.ErrorResponse "Internal Server Error"
// @Router /search/products [get]
func (h *ProductHandler) SearchProductsHandler(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")

	limit := helper.ParseLimit(limitStr)

	query, err := helper.CheckParam(r.URL.Query().Get("query"))
	if err != nil {
		HandleError(w, err, "Missing search query parameter")
		return
	}

	products, err := h.ProductRepo.SearchProducts(query, limit)
	if err != nil {
		HandleError(w, err, "Error fetching search results")
		return
	}

	products = h.ImageRepo.GetPreSignedURLs(products)

	HandleSuccessResponse(w, http.StatusOK, products)
}

func (h *ProductHandler) handleProductImageUpload(w http.ResponseWriter, r *http.Request, product *model.Product) (string, error) {
	imageData, format, err := helper.ParseProductImage(r)
	if err != nil {
		HandleError(w, err, "Error reading image")
		return "", err
	}

	s3ImageKey, err := h.ImageRepo.UploadImage(product.ProductID, r.FormValue("userId"), imageData.Bytes(), format)
	if err != nil {
		HandleError(w, err, "Error uploading image to S3")
		return "", err
	}
	return s3ImageKey, nil
}
