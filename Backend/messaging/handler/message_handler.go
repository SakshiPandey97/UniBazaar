package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"messaging/models"
	"messaging/repository"
	ws "messaging/websocket"

	"github.com/google/uuid"
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
		log.Println("Invalid user_id:", err)
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
	msg.ID = uuid.New().String()
	if err := h.repo.SaveMessage(msg); err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	h.ws.Broadcast <- msg

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "message sent"})
}

func (h *MessageHandler) GetConversationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user1IDStr := vars["user1ID"]
	user2IDStr := vars["user2ID"]

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conversations)
}

func (h *MessageHandler) GetUnreadSendersHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Error parsing user_id query parameter '%s': %v", userIDStr, err)
		http.Error(w, "Invalid or missing user_id query parameter", http.StatusBadRequest)
		return
	}
	userID := uint(userIDInt)

	senderIDs, err := h.repo.GetUnreadSenderIDs(userID)
	if err != nil {
		log.Printf("Error getting unread sender IDs from repository for user %d: %v", userID, err)
		http.Error(w, "Failed to retrieve unread sender information", http.StatusInternalServerError)
		return
	}

	if senderIDs == nil {
		senderIDs = []uint{}
	}

	responsePayload := map[string][]uint{"senderIds": senderIDs}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(responsePayload); err != nil {
		log.Printf("Error encoding JSON response for unread senders (user %d): %v", userID, err)
	}
}
