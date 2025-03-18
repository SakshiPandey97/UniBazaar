package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
		log.Printf("Error loading .env file: %v, using default settings", err)
	}

	mongoDBClient, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoDBClient.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	repo, err := repository.NewMongoProductRepository()
	if err != nil {
		log.Fatalf("Failed to create product repository: %v", err)
	}
	s3 := repository.NewS3ImageRepository()
	productHandler := handler.NewProductHandler(repo, s3)

	router := mux.NewRouter()
	router.HandleFunc("/health", healthCheck).Methods(http.MethodGet)
	routes.RegisterProductRoutes(router, productHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server running on port %s", port)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: routes.SetupCORS(router),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped")

}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Health OK"))
}
