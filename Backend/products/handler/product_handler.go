package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"web-service/model"
	"web-service/repository"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create a new product.")

	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product.ProductID = uuid.NewString()
	log.Printf("Generated Product ID: %s\n", product.ProductID)

	if err := repository.CreateProduct(product); err != nil {
		log.Printf("Error creating product: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := product.Validate(); err != nil {
		log.Printf("Validation error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
	id := mux.Vars(r)["id"]
	log.Printf("Received request to delete product with ID: %s\n", id)

	if err := repository.DeleteProduct(id); err != nil {
		log.Printf("Error deleting product with ID %s: %v\n", id, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("Product with ID %s deleted successfully.\n", id)
	w.WriteHeader(http.StatusNoContent)
}
