package main

import (
	"context"
	"fmt"
	"log"
	"messaging/db"
	"messaging/handler"
	"messaging/repository"
	"messaging/websocket"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading .env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "*"
	}

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

	r := mux.NewRouter()

	r.HandleFunc("/ws", msgHandler.HandleWebSocket)
	r.HandleFunc("/api/conversation/{user1ID}/{user2ID}", msgHandler.GetConversationHandler).Methods(http.MethodGet)
	r.HandleFunc("/messages", msgHandler.HandleSendMessage).Methods(http.MethodPost)
	r.HandleFunc("/users", userHandler.GetUsersHandler).Methods("GET")
	r.HandleFunc("/api/users/sync", userHandler.SyncUserHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/unread-senders", msgHandler.GetUnreadSendersHandler).Methods(http.MethodGet)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(r)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	go func() {
		fmt.Printf("Server started on :%s\n", port)
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
