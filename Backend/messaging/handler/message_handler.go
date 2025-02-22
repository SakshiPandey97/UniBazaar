package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"messaging/models"
	"messaging/repository"
	ws "messaging/websocket"

	gorilla "github.com/gorilla/websocket"
)

var upgrader = gorilla.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type MessageHandler struct {
	repo *repository.MessageRepository
	ws   *ws.WebSocketManager
}

func NewMessageHandler(repo *repository.MessageRepository, ws *ws.WebSocketManager) *MessageHandler {
	return &MessageHandler{repo: repo, ws: ws}
}

// WebSocket connection
func (h *MessageHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

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

	for {
		var msg models.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading JSON:", err)
			break
		}

		msg.Timestamp = time.Now().Unix()
		h.ws.SendMessage(msg)
	}
}
