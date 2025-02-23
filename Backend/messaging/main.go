package main

import (
	"fmt"
	"log"
	"net/http"

	"messaging/db"
	"messaging/handler"
	"messaging/repository"
	"messaging/websocket"
)

func main() {
	// Connect to the database
	database := db.ConnectDB()
	defer database.Close()

	// Initialize repository and WebSocket manager
	msgRepo := repository.NewMessageRepository(database)
	wsManager := websocket.NewWebSocketManager()

	// Run the WebSocket manager
	go wsManager.Run()

	// Initialize the message handler with the repository and WebSocket manager
	msgHandler := handler.NewMessageHandler(msgRepo, wsManager)

	// Define the HTTP routes
	http.HandleFunc("/ws", msgHandler.HandleWebSocket)         // WebSocket handler
	http.HandleFunc("/messages", msgHandler.HandleGetMessages) // Get latest messages

	// Start the server
	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
