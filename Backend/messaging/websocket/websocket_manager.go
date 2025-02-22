package websocket

import (
	"log"
	"sync"

	"messaging/models"
	"messaging/repository"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	UserID   uint
	SendChan chan models.Message
}

type WebSocketManager struct {
	clients    map[uint]*Client
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
	repo       *repository.MessageRepository
}

func NewWebSocketManager(repo *repository.MessageRepository) *WebSocketManager {
	return &WebSocketManager{
		clients:    make(map[uint]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		repo:       repo,
	}
}

func (wm *WebSocketManager) Run() {
	for {
		select {
		case client := <-wm.Register:
			wm.mu.Lock()
			wm.clients[client.UserID] = client
			wm.mu.Unlock()
			log.Printf("Client %d registered", client.UserID)

		case client := <-wm.Unregister:
			wm.mu.Lock()
			if _, ok := wm.clients[client.UserID]; ok {
				delete(wm.clients, client.UserID)
				close(client.SendChan)
				log.Printf("Client %d unregistered", client.UserID)

			}
			wm.mu.Unlock()
		}
	}
}

// SendMessage sends a message to the appropriate client or saves it for later delivery
func (wm *WebSocketManager) SendMessage(msg models.Message) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if client, exists := wm.clients[msg.ReceiverID]; exists {
		select {
		case client.SendChan <- msg:
			log.Printf("Message sent to client %d: %s", msg.ReceiverID, msg.Content)
		default:
			// Handle case when the channel is full or client can't receive
			log.Printf("Client %d's SendChan is full, message not sent immediately", msg.ReceiverID)
		}
	} else {
		// If client is not connected, save the message
		if err := wm.repo.SaveMessage(&msg); err != nil {
			log.Println("Error saving message:", err)
		}
	}
}
