package repository

import (
	"database/sql"
	"log"
	"messaging/models"
)

type MessageRepository struct {
	DB *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

// Save a new message
func (repo *MessageRepository) SaveMessage(msg models.Message) error {
	_, err := repo.DB.Exec(`
		INSERT INTO messages (sender_id, receiver_id, content, timestamp, read)
		VALUES ($1, $2, $3, $4, $5)`,
		msg.SenderID, msg.ReceiverID, msg.Content, msg.Timestamp, false)

	if err != nil {
		log.Println("Error saving message:", err)
		return err
	}
	return nil
}

// Get latest N messages (Fixed Query)
func (repo *MessageRepository) GetLatestMessages(limit int) ([]models.Message, error) {
	rows, err := repo.DB.Query("SELECT id, sender_id, receiver_id, content, timestamp, read FROM messages ORDER BY timestamp DESC LIMIT $1", limit)
	if err != nil {
		log.Println("Error fetching messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp, &msg.Read); err != nil {
			log.Println("Error scanning message:", err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

// Mark a message as read
func (repo *MessageRepository) MarkMessageAsRead(messageID int) error {
	_, err := repo.DB.Exec("UPDATE messages SET read = TRUE WHERE id = $1", messageID)
	if err != nil {
		log.Println("Error marking message as read:", err)
	}
	return err
}

// Get unread messages for a user
func (repo *MessageRepository) GetUnreadMessages(userID uint) ([]models.Message, error) {
	rows, err := repo.DB.Query("SELECT id, sender_id, receiver_id, content, timestamp FROM messages WHERE receiver_id = $1 AND read = FALSE", userID)
	if err != nil {
		log.Println("Error fetching unread messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp); err != nil {
			log.Println("Error scanning unread message:", err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}
