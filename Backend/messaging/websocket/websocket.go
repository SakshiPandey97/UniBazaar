package websocket

import (
	"fmt"
	"log"
	"sync"
	"time"

	"messaging/models"
	"messaging/repository"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	UserID   uint
	SendChan chan models.Message
	Manager  *WebSocketManager
}

type WebSocketManager struct {
	Clients    map[uint]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan models.Message
	Repo       *repository.MessageRepository
	mu         sync.RWMutex
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

			go client.WritePump()
			go client.ReadPump()

			go ws.SendOfflineMessages(client.UserID)

		case client := <-ws.Unregister:
			ws.mu.Lock()
			if _, ok := ws.Clients[client.UserID]; ok {
				delete(ws.Clients, client.UserID)
				close(client.SendChan)
				fmt.Printf("Client %d unregistered\n", client.UserID)
			}
			ws.mu.Unlock()

		case msg := <-ws.Broadcast:
			ws.mu.RLock()
			receiverClient, receiverExists := ws.Clients[msg.ReceiverID]
			senderClient, senderExists := ws.Clients[msg.SenderID]
			ws.mu.RUnlock()

			if receiverExists {
				select {
				case receiverClient.SendChan <- msg:
				default:
					ws.Unregister <- receiverClient
				}
			}
			if senderExists {
				select {
				case senderClient.SendChan <- msg:
				default:
					ws.Unregister <- senderClient
				}
			}
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Manager.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg models.Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		// Generate a new UUID if one doesn't exist
		if msg.ID == "" {
			msg.ID = uuid.New().String()
		}

		msg.SenderID = c.UserID
		msg.Timestamp = time.Now().Unix()
		msg.Read = false

		if err := c.Manager.Repo.SaveMessage(msg); err != nil {
			log.Println("Failed to save message(websocket):", err)
			continue
		}
		c.Manager.Broadcast <- msg
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.SendChan:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error writing to WebSocket: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (ws *WebSocketManager) SendOfflineMessages(userID uint) {
	// Get unread messages for the user.
	unreadMessages, err := ws.Repo.GetUnreadMessages(userID)
	if err != nil {
		log.Println("Error fetching unread messages:", err)
		return
	}

	ws.mu.RLock()
	client, exists := ws.Clients[userID]
	ws.mu.RUnlock()

	if exists {
		// Send each unread message to the client.
		for _, msg := range unreadMessages {
			select {
			case client.SendChan <- msg:
				//Mark the messages as read
				err := ws.Repo.MarkMessageAsRead(msg.ID)
				if err != nil {
					log.Println("Error marking message as read:", err)
				}
			case <-time.After(1 * time.Second):
				log.Println("Timeout sending offline message")
				return
			}
		}
	}
}
