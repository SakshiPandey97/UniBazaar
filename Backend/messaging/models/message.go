package models

type Message struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
	Timestamp  int64  `json:"timestamp"`
	Read       bool   `json:"read" gorm:"default:false"`
}
