package repository

import (
	"messaging/models"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) SaveMessage(msg *models.Message) error {
	return r.db.Create(msg).Error
}

func (r *MessageRepository) GetMessages(senderID, receiverID uint) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		senderID, receiverID, receiverID, senderID).Order("timestamp asc").Find(&messages).Error
	return messages, err
}
