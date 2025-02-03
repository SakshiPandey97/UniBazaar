package main

import (
	"context"
	"log"
	"net/http"

	"web-service/config"
	"web-service/repository"
	"web-service/routes"

	"github.com/gorilla/mux"
)

func main() {
	client := config.ConnectDB()
	defer client.Disconnect(context.Background())

	repository.InitProductRepo()

	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Health OK"))
	}).Methods(http.MethodGet)

	routes.RegisterProductRoutes(router)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", router)
}
