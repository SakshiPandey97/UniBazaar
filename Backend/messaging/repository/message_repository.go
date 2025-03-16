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

func (repo *MessageRepository) MarkMessageAsRead(messageID int) error {
	_, err := repo.DB.Exec("UPDATE messages SET read = TRUE WHERE id = $1", messageID)
	if err != nil {
		log.Println("Error marking message as read:", err)
	}
	return err
}

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
func (r *MessageRepository) GetConversation(user1ID uint, user2ID uint) ([]models.Conversation, error) {
	var conversations []models.Conversation

	rows, err := r.DB.Query(`
		SELECT
		    m.id,
		    m.sender_id,
		    m.receiver_id,
		    m.content,
		    m.timestamp,
			m.read,
		    s.name AS sender_name,
		    r.name AS receiver_name
		FROM
		    messages m
		JOIN
		    users s ON m.sender_id = s.id
		JOIN
		    users r ON m.receiver_id = r.id
		WHERE
		    (m.sender_id = $1 AND m.receiver_id = $2)
		    OR (m.sender_id = $2 AND m.receiver_id = $1)
		ORDER BY
		    m.timestamp ASC;
		`, user1ID, user2ID)

	if err != nil {
		log.Println("Error fetching conversation:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var conv models.Conversation
		err := rows.Scan(&conv.MessageID, &conv.SenderID, &conv.ReceiverID, &conv.MessageText, &conv.Timestamp, &conv.Read, &conv.SenderName, &conv.ReceiverName)
		if err != nil {
			log.Println("Error scanning conversation row:", err)
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error during row iteration:", err)
		return nil, err
	}

	return conversations, nil
}
