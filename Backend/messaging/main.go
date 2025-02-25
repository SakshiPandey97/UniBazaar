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
	database := db.ConnectDB()
	if database == nil {
		log.Fatal("Failed to connect to the database")
	}
	defer database.Close()

	msgRepo := repository.NewMessageRepository(database)
	wsManager := websocket.NewWebSocketManager(msgRepo)

	go wsManager.Run()

	msgHandler := handler.NewMessageHandler(msgRepo, wsManager)

	http.HandleFunc("/ws", msgHandler.HandleWebSocket)
	http.HandleFunc("/messages", msgHandler.HandleGetMessages)

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		fmt.Println("Server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v\n", err)
		}
	}()

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
