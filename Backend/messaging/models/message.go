package models

type Message struct {
	ID         int    `json:"id"`
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	Content    string `json:"content"`
	Timestamp  int64  `json:"timestamp"`
	Read       bool   `json:"read"`
}
