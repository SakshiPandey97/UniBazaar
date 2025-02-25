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

	"messaging/db"
	"messaging/handler"
	"messaging/repository"
	"messaging/websocket"
)

func main() {
	// Connect to the database
	database := db.ConnectDB()
	if database == nil {
		log.Fatal("Failed to connect to the database")
	}
	defer database.Close()

	// Initialize repository and WebSocket manager
	msgRepo := repository.NewMessageRepository(database)
	wsManager := websocket.NewWebSocketManager(msgRepo) // Pass the message repository

	// Run the WebSocket manager
	go wsManager.Run()

	// Initialize the message handler with the repository and WebSocket manager
	msgHandler := handler.NewMessageHandler(msgRepo, wsManager)

	// Define the HTTP routes
	http.HandleFunc("/ws", msgHandler.HandleWebSocket)         // WebSocket handler
	http.HandleFunc("/messages", msgHandler.HandleGetMessages) // Get latest messages

	// Start the server
	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		fmt.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v\n", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v\n", err)
	}

	fmt.Println("Server stopped gracefully")
}
