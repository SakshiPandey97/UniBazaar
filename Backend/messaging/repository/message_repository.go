package repository

import (
	"database/sql"
	"fmt"
	"log"
	"messaging/models"
)

type MessageRepository struct {
	DB *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

func (repo *MessageRepository) SaveMessage(msg models.Message) error {
	_, err := repo.DB.Exec(`
        INSERT INTO messages (id, sender_id, receiver_id, content, timestamp, read, sender_name)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		msg.ID, msg.SenderID, msg.ReceiverID, msg.Content, msg.Timestamp, msg.Read, msg.SenderName)

	if err != nil {
		log.Println("Error saving message:", err)
		return err
	}
	return nil
}

func (repo *MessageRepository) GetLatestMessages(limit int) ([]models.Message, error) {
	rows, err := repo.DB.Query("SELECT id, sender_id, receiver_id, content, timestamp, read, sender_name FROM messages ORDER BY timestamp DESC LIMIT $1", limit)
	if err != nil {
		log.Println("Error fetching messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp, &msg.Read, &msg.SenderName); err != nil {
			log.Println("Error scanning message:", err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (repo *MessageRepository) MarkMessageAsRead(messageID string) error {
	_, err := repo.DB.Exec("UPDATE messages SET read = TRUE WHERE id = $1", messageID)
	if err != nil {
		log.Println("Error marking message as read:", err)
	}
	return err
}

func (repo *MessageRepository) GetUnreadMessages(userID uint) ([]models.Message, error) {
	rows, err := repo.DB.Query("SELECT id, sender_id, receiver_id, content, timestamp, read, sender_name FROM messages WHERE receiver_id = $1 AND read = FALSE", userID)
	if err != nil {
		log.Println("Error fetching unread messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp, &msg.Read, &msg.SenderName); err != nil {
			log.Println("Error scanning unread message:", err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}
func (repo *MessageRepository) GetConversation(user1ID, user2ID uint) ([]models.Message, error) {
	rows, err := repo.DB.Query(`
        SELECT m.id, m.sender_id, m.receiver_id, m.content, m.timestamp, m.read, u.name
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        WHERE (m.sender_id = $1 AND m.receiver_id = $2) OR (m.sender_id = $2 AND m.receiver_id = $1)
        ORDER BY m.timestamp ASC`, user1ID, user2ID)
	if err != nil {
		return nil, fmt.Errorf("error querying messages: %v", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp, &msg.Read, &msg.SenderName); err != nil {
			return nil, fmt.Errorf("error scanning message row: %v", err)
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over message rows: %v", err)
	}

	return messages, nil
}

func (repo *MessageRepository) GetUnreadSenderIDs(receiverID uint) ([]uint, error) {
	query := `
        SELECT DISTINCT sender_id
        FROM messages
        WHERE receiver_id = $1 AND read = FALSE
    `
	rows, err := repo.DB.Query(query, receiverID)
	if err != nil {
		log.Printf("Error fetching unread sender IDs for receiver %d: %v", receiverID, err)
		return nil, fmt.Errorf("failed to query unread sender IDs: %w", err)
	}
	defer rows.Close()

	var senderIDs []uint
	for rows.Next() {
		var senderID uint
		if err := rows.Scan(&senderID); err != nil {
			log.Printf("Error scanning unread sender ID for receiver %d: %v", receiverID, err)
			return nil, fmt.Errorf("failed to scan sender ID: %w", err)
		}
		senderIDs = append(senderIDs, senderID)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over unread sender ID rows for receiver %d: %v", receiverID, err)
		return nil, fmt.Errorf("error during row iteration for sender IDs: %w", err)
	}

	return senderIDs, nil
}
