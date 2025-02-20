package main

import (
	"context"
	"log"
	"net/http"

	"web-service/config"
	"web-service/repository"
	"web-service/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// @title UniBazaar Products API
// @version 1.0
// @description API for managing products in the UniBazaar marketplace for university students.
// @host localhost:8080
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default settings")
	}
	mongoDBClient := config.ConnectDB()
	defer mongoDBClient.Disconnect(context.Background())

	repository.InitProductRepo()

	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health OK"))
	}).Methods(http.MethodGet)

	routes.RegisterProductRoutes(router)

	log.Println("Server running on port 8080")
	http.ListenAndServe("127.0.0.1:8080", routes.SetupCORS(router))
}
