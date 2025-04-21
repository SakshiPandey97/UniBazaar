package main

import (
	"database/sql"
	"errors"
	"messaging/models"
	"messaging/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageRepository struct {
	mock.Mock
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Close() error {
	return m.Called().Error(0)
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

func (m *MockMessageRepository) MarkMessageAsRead(messageID string) error {
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

func (m *MockMessageRepository) GetConversation(user1ID, user2ID uint) ([]models.Message, error) {
	args := m.Called(user1ID, user2ID)
	if args.Get(0) != nil {
		return args.Get(0).([]models.Message), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockMessageRepository) GetUnreadSenderIDs(receiverID uint) ([]uint, error) {
	args := m.Called(receiverID)
	if result := args.Get(0); result != nil {
		return result.([]uint), args.Error(1)
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

func TestSaveMessageError(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	msg := models.Message{
		SenderID:   1,
		ReceiverID: 2,
		Content:    "Hello, World!",
	}
	expectedError := errors.New("database error")
	mockRepo.On("SaveMessage", msg).Return(expectedError)
	err := mockRepo.SaveMessage(msg)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestGetLatestMessages(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	expectedMessages := []models.Message{
		{ID: "1", SenderID: 1, ReceiverID: 2, Content: "Test message", Timestamp: 1234567890, Read: false},
		{ID: "2", SenderID: 2, ReceiverID: 1, Content: "Another message", Timestamp: 1234567891, Read: true},
	}
	mockRepo.On("GetLatestMessages", 10).Return(expectedMessages, nil)
	messages, err := mockRepo.GetLatestMessages(10)
	assert.Nil(t, err)
	assert.Equal(t, expectedMessages, messages)
	mockRepo.AssertExpectations(t)
}

func TestGetLatestMessagesError(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	expectedError := errors.New("database error")
	mockRepo.On("GetLatestMessages", 10).Return(nil, expectedError)
	messages, err := mockRepo.GetLatestMessages(10)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, messages)
	mockRepo.AssertExpectations(t)
}

func TestMarkMessageAsRead(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	mockRepo.On("MarkMessageAsRead", "1").Return(nil)
	err := mockRepo.MarkMessageAsRead("1")
	assert.Nil(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMarkMessageAsReadError(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	expectedError := errors.New("database error")
	mockRepo.On("MarkMessageAsRead", "1").Return(expectedError)
	err := mockRepo.MarkMessageAsRead("1")
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUnreadMessages(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	expectedMessages := []models.Message{
		{ID: "1", SenderID: 1, ReceiverID: 2, Content: "Unread message", Timestamp: 1234567892, Read: false},
	}
	mockRepo.On("GetUnreadMessages", uint(2)).Return(expectedMessages, nil)
	messages, err := mockRepo.GetUnreadMessages(2)
	assert.Nil(t, err)
	assert.Equal(t, expectedMessages, messages)
	mockRepo.AssertExpectations(t)
}

func TestGetUnreadMessagesError(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	expectedError := errors.New("database error")
	mockRepo.On("GetUnreadMessages", uint(2)).Return(nil, expectedError)
	messages, err := mockRepo.GetUnreadMessages(2)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, messages)
	mockRepo.AssertExpectations(t)
}

func TestGetConversation(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	expectedMessages := []models.Message{
		{ID: "1", SenderID: 1, ReceiverID: 2, Content: "Message 1", Timestamp: 1234567890, Read: false},
		{ID: "2", SenderID: 2, ReceiverID: 1, Content: "Message 2", Timestamp: 1234567891, Read: true},
	}
	mockRepo.On("GetConversation", uint(1), uint(2)).Return(expectedMessages, nil)
	messages, err := mockRepo.GetConversation(1, 2)
	assert.Nil(t, err)
	assert.Equal(t, expectedMessages, messages)
	mockRepo.AssertExpectations(t)
}

func TestGetConversationError(t *testing.T) {
	mockRepo := new(MockMessageRepository)
	expectedError := errors.New("database error")
	mockRepo.On("GetConversation", uint(1), uint(2)).Return(nil, expectedError)
	messages, err := mockRepo.GetConversation(1, 2)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, messages)
	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	expectedUsers := []models.User{
		{ID: 1, Name: "User 1", Email: "user1@example.com"},
		{ID: 2, Name: "User 2", Email: "user2@example.com"},
	}
	mockRepo.On("GetAllUsers").Return(expectedUsers, nil)
	users, err := mockRepo.GetAllUsers()
	assert.Nil(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}

func TestGetAllUsersError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	expectedError := errors.New("database error")
	mockRepo.On("GetAllUsers").Return(nil, expectedError)
	users, err := mockRepo.GetAllUsers()
	assert.Equal(t, expectedError, err)
	assert.Nil(t, users)
	mockRepo.AssertExpectations(t)
}

func TestNewUserRepository(t *testing.T) {
	db := &sql.DB{}

	userRepo := repository.NewUserRepository(db)

	assert.NotNil(t, userRepo)

	assert.Equal(t, db, userRepo.DB)
}

func TestNewMessageRepository(t *testing.T) {
	db := &sql.DB{}

	messageRepo := repository.NewMessageRepository(db)

	assert.NotNil(t, messageRepo)

	assert.Equal(t, db, messageRepo.DB)
}

// Add these tests to msg_unit_test.go

func TestGetUnreadSenderIDs(t *testing.T) {
	// Arrange
	mockRepo := new(MockMessageRepository)
	receiverID := uint(10)
	expectedSenderIDs := []uint{1, 5, 3} // Example sender IDs

	// Setup mock expectation
	mockRepo.On("GetUnreadSenderIDs", receiverID).Return(expectedSenderIDs, nil)

	// Act
	senderIDs, err := mockRepo.GetUnreadSenderIDs(receiverID)

	// Assert
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, expectedSenderIDs, senderIDs, "Expected sender IDs to match")
	mockRepo.AssertExpectations(t) // Verify that the mock was called as expected
}

func TestGetUnreadSenderIDsError(t *testing.T) {
	// Arrange
	mockRepo := new(MockMessageRepository)
	receiverID := uint(10)
	expectedError := errors.New("database query failed")

	// Setup mock expectation to return an error
	mockRepo.On("GetUnreadSenderIDs", receiverID).Return(nil, expectedError)

	// Act
	senderIDs, err := mockRepo.GetUnreadSenderIDs(receiverID)

	// Assert
	assert.Error(t, err, "Expected an error")
	assert.Equal(t, expectedError, err, "Expected the specific database error")
	assert.Nil(t, senderIDs, "Expected sender IDs slice to be nil on error")
	mockRepo.AssertExpectations(t) // Verify that the mock was called as expected
}
