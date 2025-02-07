package routes

import (
	"web-service/handler"

	"github.com/gorilla/mux"
)

func RegisterProductRoutes(router *mux.Router) {
	router.HandleFunc("/products", handler.CreateProductHandler).Methods("POST")
	router.HandleFunc("/products", handler.GetAllProductsHandler).Methods("GET")
	router.HandleFunc("/products/{id}", handler.GetProductByIDHandler).Methods("GET")
	router.HandleFunc("/products/{id}", handler.UpdateProductHandler).Methods("PUT")
	router.HandleFunc("/products/{id}", handler.DeleteProductHandler).Methods("DELETE")
}
