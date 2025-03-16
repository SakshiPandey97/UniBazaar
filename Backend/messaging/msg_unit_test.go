package main

import (
	"testing"
	"messaging/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) SaveMessage(msg models.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockMessageRepository) GetLatestMessages(limit int) ([]models.Message, error) {
	args := m.Called(limit)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Message), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMessageRepository) MarkMessageAsRead(messageID int) error {
	args := m.Called(messageID)
	return args.Error(0)
}

func (m *MockMessageRepository) GetUnreadMessages(userID uint) ([]models.Message, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Message), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMessageRepository) GetConversation(userID uint) ([]models.Message, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Message), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestSaveMessage(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	msg := models.Message{
		SenderID:   1,
		ReceiverID: 2,
		Content:    "Hello, World!",
	}
	mockRepo.On("SaveMessage", msg).Return(nil)
	err := mockRepo.SaveMessage(msg)
	assert.Nil(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetLatestMessages(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	expectedMessages := []models.Message{
		{ID: 1, SenderID: 1, ReceiverID: 2, Content: "Test message", Timestamp: 1234567890, Read: false},
		{ID: 2, SenderID: 2, ReceiverID: 1, Content: "Another message", Timestamp: 1234567891, Read: true},
	}
	mockRepo.On("GetLatestMessages", 10).Return(expectedMessages, nil)
	messages, err := mockRepo.GetLatestMessages(10)
	assert.Nil(t, err)
	assert.Equal(t, expectedMessages, messages)
	mockRepo.AssertExpectations(t)
}

func TestMarkMessageAsRead(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	mockRepo.On("MarkMessageAsRead", 1).Return(nil)
	err := mockRepo.MarkMessageAsRead(1)
	assert.Nil(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUnreadMessages(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	expectedMessages := []models.Message{
		{ID: 1, SenderID: 1, ReceiverID: 2, Content: "Unread message", Timestamp: 1234567892, Read: false},
	}
	mockRepo.On("GetUnreadMessages", uint(2)).Return(expectedMessages, nil)
	messages, err := mockRepo.GetUnreadMessages(2)
	assert.Nil(t, err)
	assert.Equal(t, expectedMessages, messages)
	mockRepo.AssertExpectations(t)
}

func TestGetConversation(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	expectedMessages := []models.Message{
		{ID: 1, SenderID: 1, ReceiverID: 2, Content: "Message 1", Timestamp: 1234567890, Read: false},
		{ID: 2, SenderID: 2, ReceiverID: 1, Content: "Message 2", Timestamp: 1234567891, Read: true},
	}
	mockRepo.On("GetConversation", uint(1)).Return(expectedMessages, nil)
	messages, err := mockRepo.GetConversation(1)
	assert.Nil(t, err)
	assert.Equal(t, expectedMessages, messages)
	mockRepo.AssertExpectations(t)
}