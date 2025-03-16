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

	"github.com/gorilla/mux"
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

func (h *MessageHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		conn.Close()
		return
	}

	client := &ws.Client{
		Conn:     conn,
		UserID:   uint(userID),
		SendChan: make(chan models.Message, 256),
		Manager:  h.ws,
	}

	h.ws.Register <- client

}

func (h *MessageHandler) HandleSendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	msg.Timestamp = time.Now().Unix()
	msg.Read = false

	h.ws.Broadcast <- msg

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "message sent"})
}

// func (h *MessageHandler) HandleGetMessages(w http.ResponseWriter, r *http.Request) {
// 	senderIDStr := r.URL.Query().Get("sender_id")
// 	receiverIDStr := r.URL.Query().Get("receiver_id")

// 	if senderIDStr == "" || receiverIDStr == "" {
// 		http.Error(w, "sender_id and receiver_id are required", http.StatusBadRequest)
// 		return
// 	}

// 	senderID, err := strconv.Atoi(senderIDStr)
// 	if err != nil {
// 		http.Error(w, "Invalid sender_id", http.StatusBadRequest)
// 		return
// 	}

// 	if err != nil {
// 		http.Error(w, "Invalid receiver_id", http.StatusBadRequest)
// 		return
// 	}

// 	messages, err := h.repo.GetConversation(uint(senderID))
// 	if err != nil {
// 		fmt.Println("Database error:", err)
// 		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
// 		return
// 	}

//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusOK)
//		if err := json.NewEncoder(w).Encode(messages); err != nil {
//			http.Error(w, "Error encoding response", http.StatusInternalServerError)
//		}
//	}
func (h *MessageHandler) GetConversationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user1IDStr := vars["user1ID"]
	user2IDStr := vars["user2ID"]

	// Basic validation and type conversion
	user1ID, err := strconv.ParseUint(user1IDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user1ID", http.StatusBadRequest)
		return
	}
	user2ID, err := strconv.ParseUint(user2IDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid user2ID", http.StatusBadRequest)
		return
	}

	conversations, err := h.repo.GetConversation(uint(user1ID), uint(user2ID))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting conversation: %v", err), http.StatusInternalServerError)
		return
	}

	// Respond with the messages
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}
