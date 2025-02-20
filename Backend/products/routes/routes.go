package routes

import (
	"net/http"
	"web-service/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func RegisterProductRoutes(router *mux.Router) {
	// Define product routes
	router.HandleFunc("/products", handler.CreateProductHandler).Methods("POST")
	router.HandleFunc("/products", handler.GetAllProductsHandler).Methods("GET")
	router.HandleFunc("/products/{UserId}", handler.GetAllProductsByUserIDHandler).Methods("GET")
	router.HandleFunc("/products/{UserId}/{ProductId}", handler.UpdateProductHandler).Methods("PUT")
	router.HandleFunc("/products/{UserId}/{ProductId}", handler.DeleteProductHandler).Methods("DELETE")
}

func SetupCORS(router *mux.Router) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins (change for security)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(router)
}
