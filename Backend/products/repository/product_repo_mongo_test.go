package repository

import (
	"errors"
	"testing"
	"web-service/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProductRepository is the mock of the ProductRepository interface
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) CreateProduct(product model.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetAllProducts() ([]model.Product, error) {
	args := m.Called()
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockProductRepository) GetProductsByUserID(userID int) ([]model.Product, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateProduct(userID int, productID string, product model.Product) error {
	args := m.Called(userID, productID, product)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteProduct(userID int, productID string) error {
	args := m.Called(userID, productID)
	return args.Error(0)
}

func (m *MockProductRepository) FindProductByUserAndId(userID int, productID string) (*model.Product, error) {
	args := m.Called(userID, productID)
	return args.Get(0).(*model.Product), args.Error(1)
}

func TestCreateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)

	product := model.Product{UserID: 1, ProductID: "prod123"}

	mockRepo.On("CreateProduct", product).Return(nil)

	err := mockRepo.CreateProduct(product)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_Error(t *testing.T) {
	mockRepo := new(MockProductRepository)

	product := model.Product{UserID: 1, ProductID: "prod123"}

	mockRepo.On("CreateProduct", product).Return(errors.New("insert error"))

	err := mockRepo.CreateProduct(product)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAllProducts(t *testing.T) {
	mockRepo := new(MockProductRepository)

	mockRepo.On("GetAllProducts").Return([]model.Product{
		{UserID: 1, ProductID: "prod123"},
		{UserID: 2, ProductID: "prod456"},
	}, nil)

	products, err := mockRepo.GetAllProducts()

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "prod123", products[0].ProductID)
	mockRepo.AssertExpectations(t)
}

func TestGetProductsByUserID(t *testing.T) {
	mockRepo := new(MockProductRepository)
	userID := 1

	mockRepo.On("GetProductsByUserID", userID).Return([]model.Product{
		{UserID: userID, ProductID: "prod123"},
		{UserID: userID, ProductID: "prod456"},
	}, nil)

	products, err := mockRepo.GetProductsByUserID(userID)

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "prod123", products[0].ProductID)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)

	userID := 1
	productID := "prod123"
	product := model.Product{UserID: userID, ProductID: productID}

	mockRepo.On("UpdateProduct", userID, productID, product).Return(nil)

	err := mockRepo.UpdateProduct(userID, productID, product)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)

	userID := 1
	productID := "prod123"

	mockRepo.On("DeleteProduct", userID, productID).Return(nil)

	err := mockRepo.DeleteProduct(userID, productID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestFindProductByUserAndId(t *testing.T) {
	mockRepo := new(MockProductRepository)

	userID := 1
	productID := "prod123"

	mockRepo.On("FindProductByUserAndId", userID, productID).Return(&model.Product{
		UserID:    userID,
		ProductID: productID,
	}, nil)

	product, err := mockRepo.FindProductByUserAndId(userID, productID)

	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, productID, product.ProductID)
	mockRepo.AssertExpectations(t)
}
