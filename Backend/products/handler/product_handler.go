package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"web-service/model"
	"web-service/repository"

	"github.com/gorilla/mux"
)

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to create a new product.")

	userProduct := model.UserProduct{}

	err := repository.CreateProduct(userProduct, r)
	if err != nil {
		log.Printf("Error creating product: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Product created successfully\n")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userProduct.Products[0])
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
	userId := r.URL.Query().Get("userId")
	productId := r.URL.Query().Get("productId")

	if userId == "" || productId == "" {
		http.Error(w, "Missing userId or productId in query parameters", http.StatusBadRequest)
		return
	}

	err := repository.UpdateProduct(w, r, userId, productId)
	if err != nil {
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
