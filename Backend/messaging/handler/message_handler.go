package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"messaging/models"
	"messaging/repository"
	ws "messaging/websocket"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type MessageHandler struct {
	repo *repository.MessageRepository
	ws   *ws.WebSocketManager
}

func NewMessageHandler(repo *repository.MessageRepository, ws *ws.WebSocketManager) *MessageHandler {
	return &MessageHandler{repo: repo, ws: ws}
}

// WebSocket connection handler
func (h *MessageHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	// Get userID from query params
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		conn.Close()
		return
	}

	client := &ws.Client{
		Conn:     conn,
		UserID:   uint(userID),
		SendChan: make(chan models.Message),
	}

	h.ws.Register <- client

	// Ensure client is unregistered and resources are cleaned up
	defer func() {
		h.ws.Unregister <- client
		close(client.SendChan) // Close channel to prevent goroutine leaks
		conn.Close()
	}()

	// Dedicated goroutine to send messages to the client
	go func() {
		for msg := range client.SendChan {
			if err := conn.WriteJSON(msg); err != nil {
				fmt.Println("Error writing JSON:", err)
				break
			}
		}
	}()

	// Read incoming messages from the WebSocket
	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("WebSocket read error:", err)
			break
		}

		// Set the timestamp before saving
		msg.Timestamp = time.Now().Unix()

		// Save message to repository
		if err := h.repo.SaveMessage(msg); err != nil {
			fmt.Println("Error saving message:", err)
			continue // Skip sending if saving fails
		}

		// Send message to receiver
		h.ws.SendOfflineMessages(client.UserID)
	}
}

// HTTP handler to get recent messages
func (h *MessageHandler) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		limit = "15"
	}

	numMessages, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "Invalid limit value", http.StatusBadRequest)
		return
	}

	messages, err := h.repo.GetLatestMessages(numMessages)
	if err != nil {
		fmt.Println("Database error:", err)
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
