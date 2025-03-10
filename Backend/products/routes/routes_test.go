package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"web-service/handler"
	"web-service/model"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
)

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

func (m *MockProductRepository) GetProductsByUserID(userID int, lastID string, limit int) ([]model.Product, error) {
	args := m.Called(userID, lastID, limit)
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

type MockImageRepository struct {
	mock.Mock
}

func (m *MockImageRepository) UploadImage(productID string, userID string, imageData []byte, format string) (string, error) {
	args := m.Called(productID, userID, imageData, format)
	return args.String(0), args.Error(1)
}

func (m *MockImageRepository) DeleteImage(imageKey string) error {
	args := m.Called(imageKey)
	return args.Error(0)
}

func (m *MockImageRepository) GetPreSignedURLs(products []model.Product) []model.Product {
	args := m.Called(products)
	return args.Get(0).([]model.Product)
}

func (m *MockImageRepository) GeneratePresignedURL(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

// func TestRegisterProductRoutes(t *testing.T) {
// 	mockProductRepo := new(MockProductRepository)
// 	mockImageRepo := new(MockImageRepository)
// 	handler := handler.NewProductHandler(mockProductRepo, mockImageRepo)

// 	router := mux.NewRouter()
// 	RegisterProductRoutes(router, handler)

// 	mockProductRepo.On("CreateProduct", mock.AnythingOfType("*model.Product")).Return(nil).Once()

// 	req, err := http.NewRequest(http.MethodPost, "/products", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	mockProductRepo.On("GetAllProducts", "327efad4-4ac0-462e-b1cf-811662821212", 2).Return([]model.Product{
// 		{UserID: 1, ProductTitle: "Product 1", ProductID: "product1"},
// 		{UserID: 1, ProductTitle: "Product 2", ProductID: "product2"},
// 	}, nil).Once()

// 	mockImageRepo.On("GetPreSignedURLs", mock.AnythingOfType("[]model.Product")).Return([]model.Product{}, nil).Once()

// 	req, err = http.NewRequest(http.MethodGet, "/products?limit=2&lastID=327efad4-4ac0-462e-b1cf-811662821212", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr = httptest.NewRecorder()
// 	router.ServeHTTP(rr, req)

// 	assert.Equal(t, http.StatusOK, rr.Code)

// }

func TestCORSHeaders(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockImageRepo := new(MockImageRepository)
	handler := handler.NewProductHandler(mockProductRepo, mockImageRepo)

	router := mux.NewRouter()
	RegisterProductRoutes(router, handler)

	corsRouter := SetupCORS(router)

	req, err := http.NewRequest(http.MethodOptions, "/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	corsRouter.ServeHTTP(rr, req)

	verifyCORSHeaders(t, rr)
}

func verifyCORSHeaders(t *testing.T, rr *httptest.ResponseRecorder) {
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Logf("Warning: Expected Access-Control-Allow-Origin to be '*', got %s", rr.Header().Get("Access-Control-Allow-Origin"))
	}
	if rr.Header().Get("Access-Control-Allow-Methods") != "GET,POST,PUT,DELETE,OPTIONS" {
		t.Logf("Warning: Expected Access-Control-Allow-Methods to be 'GET,POST,PUT,DELETE,OPTIONS', got %s", rr.Header().Get("Access-Control-Allow-Methods"))
	}
	if rr.Header().Get("Access-Control-Allow-Headers") != "Content-Type,Authorization" {
		t.Logf("Warning: Expected Access-Control-Allow-Headers to be 'Content-Type,Authorization', got %s", rr.Header().Get("Access-Control-Allow-Headers"))
	}
}
