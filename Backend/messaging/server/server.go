package server

import (
	"fmt"
	"log"
	"net/http"

	"messaging/config"
	"messaging/handler"
	"messaging/repository"
	"messaging/routes"
	"messaging/websocket"

	"github.com/gorilla/mux"
)

func StartServer() {
	config.InitDatabase()

	repo := repository.NewMessageRepository(config.DB)
	wsManager := websocket.NewWebSocketManager(repo)
	go wsManager.Run()

	messageHandler := handler.NewMessageHandler(repo, wsManager)

	router := mux.NewRouter()
	routes.SetupRoutes(router, messageHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
