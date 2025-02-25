package websocket

import (
	"fmt"
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
	Clients    map[uint]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan models.Message
	Repo       *repository.MessageRepository
	mu         sync.Mutex
}

func NewWebSocketManager(repo *repository.MessageRepository) *WebSocketManager {
	return &WebSocketManager{
		Clients:    make(map[uint]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan models.Message),
		Repo:       repo,
	}
}

func (ws *WebSocketManager) Run() {
	for {
		select {
		case client := <-ws.Register:
			ws.mu.Lock()
			ws.Clients[client.UserID] = client
			ws.mu.Unlock()
			fmt.Printf("Client %d registered\n", client.UserID)

			// Send offline messages upon reconnection
			ws.SendOfflineMessages(client.UserID)

		case client := <-ws.Unregister:
			ws.mu.Lock()
			delete(ws.Clients, client.UserID)
			close(client.SendChan)
			ws.mu.Unlock()
			fmt.Printf("Client %d unregistered\n", client.UserID)

		case msg := <-ws.Broadcast:
			ws.mu.Lock()
			receiverClient, exists := ws.Clients[msg.ReceiverID]
			ws.mu.Unlock()

			if exists {
				// Send the message to the online user
				receiverClient.SendChan <- msg

				// Mark message as read in the database
				ws.Repo.MarkMessageAsRead(msg.ID)
			} else {
				// Store the message for offline retrieval
				err := ws.Repo.SaveMessage(msg)
				if err != nil {
					log.Println("Error saving message:", err)
				}
			}
		}
	}
}

// Send messages to a reconnected user
func (ws *WebSocketManager) SendOfflineMessages(userID uint) {
	messages, err := ws.Repo.GetUnreadMessages(userID)
	if err != nil {
		log.Println("Error fetching unread messages:", err)
		return
	}

	ws.mu.Lock()
	client, exists := ws.Clients[userID]
	ws.mu.Unlock()

	if exists {
		for _, msg := range messages {
			client.SendChan <- msg
			ws.Repo.MarkMessageAsRead(msg.ID) // Mark messages as read
		}
	}
}
