package main

import (
	"testing"
	"users/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserModel struct {
	mock.Mock
}

func (m *MockUserModel) Insert(id int, name, email, password, phone string) error {
	args := m.Called(id, name, email, password, phone)
	return args.Error(0)
}

func (m *MockUserModel) Read(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserModel) UpdateName(email, newName string) error {
	args := m.Called(email, newName)
	return args.Error(0)
}

func (m *MockUserModel) UpdatePhone(email, newPhone string) error {
	args := m.Called(email, newPhone)
	return args.Error(0)
}

func (m *MockUserModel) Delete(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserModel) InitiatePasswordReset(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func (m *MockUserModel) VerifyResetCodeAndSetNewPassword(email, code, newPassword string) error {
	args := m.Called(email, code, newPassword)
	return args.Error(0)
}

func TestUserInsert(t *testing.T) {
	mockUserModel := new(MockUserModel)
	mockUserModel.On("Insert", 1, "John Doe", "john@ufl.edu", "StrongP@ssw0rd!", "5551234567").Return(nil)
	err := mockUserModel.Insert(1, "John Doe", "john@ufl.edu", "StrongP@ssw0rd!", "5551234567")
	assert.Nil(t, err)
	mockUserModel.AssertExpectations(t)
}

func TestUserRead(t *testing.T) {
	mockUserModel := new(MockUserModel)
	expectedUser := &models.User{
		UserID:   1,
		Name:     "John Doe",
		Email:    "john@ufl.edu",
		Password: "hashedpassword",
		Phone:    "5551234567",
	}
	mockUserModel.On("Read", "john@ufl.edu").Return(expectedUser, nil)
	user, err := mockUserModel.Read("john@ufl.edu")
	assert.Nil(t, err)
	assert.Equal(t, expectedUser, user)
	mockUserModel.AssertExpectations(t)
}

func TestUpdateUserName(t *testing.T) {
	mockUserModel := new(MockUserModel)
	mockUserModel.On("UpdateName", "john@ufl.edu", "Johnathan Doe").Return(nil)
	err := mockUserModel.UpdateName("john@ufl.edu", "Johnathan Doe")
	assert.Nil(t, err)
	mockUserModel.AssertExpectations(t)
}

func TestUpdateUserPhone(t *testing.T) {
	mockUserModel := new(MockUserModel)
	mockUserModel.On("UpdatePhone", "john@ufl.edu", "5559998888").Return(nil)
	err := mockUserModel.UpdatePhone("john@ufl.edu", "5559998888")
	assert.Nil(t, err)
	mockUserModel.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockUserModel := new(MockUserModel)
	mockUserModel.On("Delete", "john@ufl.edu").Return(nil)
	err := mockUserModel.Delete("john@ufl.edu")
	assert.Nil(t, err)
	mockUserModel.AssertExpectations(t)
}

func TestPasswordReset(t *testing.T) {
	mockUserModel := new(MockUserModel)
	mockUserModel.On("InitiatePasswordReset", "john@ufl.edu").Return(nil)
	err := mockUserModel.InitiatePasswordReset("john@ufl.edu")
	assert.Nil(t, err)
	mockUserModel.AssertExpectations(t)
}

func TestVerifyResetCodeAndSetNewPassword(t *testing.T) {
	mockUserModel := new(MockUserModel)
	mockUserModel.On("VerifyResetCodeAndSetNewPassword", "john@ufl.edu", "123456", "NewStrongP@ssw0rd!").Return(nil)
	err := mockUserModel.VerifyResetCodeAndSetNewPassword("john@ufl.edu", "123456", "NewStrongP@ssw0rd!")
	assert.Nil(t, err)
	mockUserModel.AssertExpectations(t)
}

func TestValidateEduEmail(t *testing.T) {
	err := models.ValidateEduEmail("user@ufl.edu")
	assert.Nil(t, err)
	err = models.ValidateEduEmail("user@gmail.com")
	assert.NotNil(t, err)
}

func TestValidatePassword(t *testing.T) {
	err := models.ValidatePassword("abc")
	assert.NotNil(t, err)
	err = models.ValidatePassword("Str0ngPassword123!")
	assert.Nil(t, err)
}

func TestValidatePhone(t *testing.T) {
	err := models.ValidatePhone("5551234567")
	assert.Nil(t, err)
	err = models.ValidatePhone("+15551234567")
	assert.Nil(t, err)
	err = models.ValidatePhone("919876543210")
	assert.NotNil(t, err)
}
