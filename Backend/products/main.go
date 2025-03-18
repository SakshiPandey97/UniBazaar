package main

import (
	"context"
	"log"
	"net/http"

	"web-service/config"
	"web-service/handler"
	"web-service/repository"
	"web-service/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// @title UniBazaar Products API
// @version 1.0
// @description API for managing products in the UniBazaar marketplace for university students.
// @host unibazaar-products.azurewebsites.net
// @schemes https
// @BasePath /
// @contact.name Avaneesh Khandekar
// @contact.email avaneesh.khandekar@gmail.com
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default settings")
	}

	mongoDBClient := config.ConnectDB()
	defer mongoDBClient.Disconnect(context.Background())

	repo := repository.NewMongoProductRepository()
	s3 := repository.NewS3ImageRepository()
	handler := handler.NewProductHandler(repo, s3)

	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health OK"))
	}).Methods(http.MethodGet)

	routes.RegisterProductRoutes(router, handler)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", routes.SetupCORS(router))
}
