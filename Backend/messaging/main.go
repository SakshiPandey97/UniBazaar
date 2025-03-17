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

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database := db.ConnectDB()
	if database == nil {
		log.Fatal("Failed to connect to the database")
	}
	defer database.Close()

	msgRepo := repository.NewMessageRepository(database)
	userRepo := repository.NewUserRepository(database)
	wsManager := websocket.NewWebSocketManager(msgRepo)

	go wsManager.Run()

	msgHandler := handler.NewMessageHandler(msgRepo, wsManager)

	userHandler := handler.NewUserHandler(userRepo)

	// Create a new ServeMux
	r := mux.NewRouter()

	// Register your handlers
	r.HandleFunc("/ws", msgHandler.HandleWebSocket)
	r.HandleFunc("/api/conversation/{user1ID}/{user2ID}", msgHandler.GetConversationHandler).Methods(http.MethodGet)
	r.HandleFunc("/messages", msgHandler.HandleSendMessage).Methods(http.MethodPost)
	r.HandleFunc("/users", userHandler.GetUsersHandler).Methods("GET")

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow requests from your frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		Debug:            true, // Log CORS-related issues
	})

	// Wrap your ServeMux with the CORS middleware
	handler := c.Handler(r)

	// Create the server with the CORS-wrapped handler
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
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
