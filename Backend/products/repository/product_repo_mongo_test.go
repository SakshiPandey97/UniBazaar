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

func (m *MockProductRepository) GetAllProducts(lastID string, limit int) ([]model.Product, error) {
	args := m.Called(lastID, limit)
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

func (m *MockProductRepository) SearchProducts(query string, limit int) ([]model.Product, error) {
	args := m.Called(query, limit)
	return args.Get(0).([]model.Product), args.Error(1)
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

	lastID := "prod100"
	limit := 2

	mockRepo.On("GetAllProducts", lastID, limit).Return([]model.Product{
		{UserID: 1, ProductID: "prod123"},
		{UserID: 2, ProductID: "prod456"},
	}, nil)

	products, err := mockRepo.GetAllProducts(lastID, limit)

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

func TestSearchProducts_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	query := "laptop"
	limit := 2

	expectedProducts := []model.Product{
		{UserID: 1, ProductID: "prod123", ProductTitle: "Laptop X1"},
		{UserID: 2, ProductID: "prod456", ProductTitle: "Gaming Laptop"},
	}

	mockRepo.On("SearchProducts", query, limit).Return(expectedProducts, nil)

	products, err := mockRepo.SearchProducts(query, limit)

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Laptop X1", products[0].ProductTitle)
	mockRepo.AssertExpectations(t)
}

func TestSearchProducts_NoResults(t *testing.T) {
	mockRepo := new(MockProductRepository)
	query := "laptop"
	limit := 2

	mockRepo.On("SearchProducts", query, limit).Return([]model.Product{}, errors.New("no products found"))

	products, err := mockRepo.SearchProducts(query, limit)

	assert.Error(t, err)
	assert.Equal(t, "no products found", err.Error())
	assert.Len(t, products, 0)
	mockRepo.AssertExpectations(t)
}

func TestSearchProducts_DatabaseError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	query := "laptop"
	limit := 2

	mockRepo.On("SearchProducts", query, limit).Return([]model.Product{}, errors.New("database error"))

	products, err := mockRepo.SearchProducts(query, limit)

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	assert.Len(t, products, 0)
	mockRepo.AssertExpectations(t)
}

func TestGetAllProducts_Pagination_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)

	lastID := "prod100"
	limit := 2

	mockRepo.On("GetAllProducts", lastID, limit).Return([]model.Product{
		{UserID: 1, ProductID: "prod123"},
		{UserID: 2, ProductID: "prod456"},
	}, nil)

	products, err := mockRepo.GetAllProducts(lastID, limit)

	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "prod123", products[0].ProductID)
	mockRepo.AssertExpectations(t)
}

func TestGetAllProducts_Pagination_NoLastID_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)

	lastID := ""
	limit := 3

	mockRepo.On("GetAllProducts", lastID, limit).Return([]model.Product{
		{UserID: 1, ProductID: "prod123"},
		{UserID: 2, ProductID: "prod456"},
		{UserID: 3, ProductID: "prod789"},
	}, nil)

	products, err := mockRepo.GetAllProducts(lastID, limit)

	assert.NoError(t, err)
	assert.Len(t, products, 3)
	assert.Equal(t, "prod123", products[0].ProductID)
	mockRepo.AssertExpectations(t)
}

func TestGetAllProducts_Pagination_EmptyResults(t *testing.T) {
	mockRepo := new(MockProductRepository)

	lastID := "prod999"
	limit := 2

	mockRepo.On("GetAllProducts", lastID, limit).Return([]model.Product{}, nil)

	products, err := mockRepo.GetAllProducts(lastID, limit)

	assert.NoError(t, err)
	assert.Len(t, products, 0)
	mockRepo.AssertExpectations(t)
}
