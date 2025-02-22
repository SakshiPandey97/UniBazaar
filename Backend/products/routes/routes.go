package routes

import (
	"net/http"
	"web-service/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func RegisterProductRoutes(router *mux.Router, productHandler *handler.ProductHandler) {
	router.HandleFunc("/products", productHandler.CreateProductHandler).Methods("POST")
	router.HandleFunc("/products", productHandler.GetAllProductsHandler).Methods("GET")
	router.HandleFunc("/products/{UserId}", productHandler.GetAllProductsByUserIDHandler).Methods("GET")
	router.HandleFunc("/products/{UserId}/{ProductId}", productHandler.UpdateProductHandler).Methods("PUT")
	router.HandleFunc("/products/{UserId}/{ProductId}", productHandler.DeleteProductHandler).Methods("DELETE")
}

func SetupCORS(router *mux.Router) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins (adjust for security)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(router)
}
