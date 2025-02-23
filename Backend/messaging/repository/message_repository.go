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

// Method to save a message to the database
func (repo *MessageRepository) SaveMessage(msg models.Message) error {
	// Insert the message into the database
	_, err := repo.DB.Exec(`
		INSERT INTO messages (sender_id, receiver_id, content, timestamp)
		VALUES ($1, $2, $3, $4)`,
		msg.SenderID, msg.ReceiverID, msg.Content, msg.Timestamp)

	if err != nil {
		log.Println("Error saving message:", err)
		return err
	}

	return nil
}

// Method to fetch the latest N messages
func (repo *MessageRepository) GetLatestMessages(limit int) ([]models.Message, error) {
	rows, err := repo.DB.Query("SELECT id, sender_id, receiver_id, content, timestamp FROM messages ORDER BY timestamp DESC LIMIT $1?", limit)
	if err != nil {
		log.Println("Error fetching messages:", err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.Timestamp); err != nil {
			log.Println("Error scanning message:", err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error in row iteration:", err)
		return nil, err
	}

	return messages, nil
}
