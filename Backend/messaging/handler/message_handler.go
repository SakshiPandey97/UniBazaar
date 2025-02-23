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
	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))

	client := &ws.Client{
		Conn:     conn,
		UserID:   uint(userID),
		SendChan: make(chan models.Message),
	}

	h.ws.Register <- client

	defer func() {
		h.ws.Unregister <- client
		conn.Close()
	}()

	go func() {
		for msg := range client.SendChan {
			conn.WriteJSON(msg)
		}
	}()

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading JSON:", err)
			break
		}

		msg.Timestamp = time.Now().Unix()

		h.repo.SaveMessage(msg)

		h.ws.SendMessageToReceiver(msg)
	}
}

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
		http.Error(w, "Error fetching messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
