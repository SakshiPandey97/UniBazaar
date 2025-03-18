package models

type Conversation struct {
	MessageID    int    `json:"message_id"`
	SenderID     uint   `json:"sender_id"`
	ReceiverID   uint   `json:"receiver_id"`
	MessageText  string `json:"message_text"`
	Timestamp    int64  `json:"timestamp"`
	SenderName   string `json:"sender_name"`
	ReceiverName string `json:"receiver_name"`
	Read         bool   `json:"read"`
}
