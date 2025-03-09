package handler

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"web-service/model"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
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

func CreateMockImage(format string) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}

	var buf bytes.Buffer
	switch format {
	case "jpeg":
		err := jpeg.Encode(&buf, img, nil)
		if err != nil {
			return nil, fmt.Errorf("error encoding JPEG: %v", err)
		}
	case "png":
		err := png.Encode(&buf, img)
		if err != nil {
			return nil, fmt.Errorf("error encoding PNG: %v", err)
		}
	case "gif":
		return nil, fmt.Errorf("unsupported format: gif")
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}

	return buf.Bytes(), nil
}

func TestCreateProductHandler(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockImageRepo := new(MockImageRepository)
	handler := NewProductHandler(mockProductRepo, mockImageRepo)

	file, err := CreateMockImage("jpeg")
	if err != nil {
		t.Fatalf("Error creating mock image: %v", err)
	}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	_ = writer.WriteField("userId", "1")
	_ = writer.WriteField("productTitle", "Test Product")
	_ = writer.WriteField("productDescription", "A sample product description")
	_ = writer.WriteField("productPostDate", "03-03-2025")
	_ = writer.WriteField("productCondition", "4")
	_ = writer.WriteField("productPrice", "9.99")
	_ = writer.WriteField("productLocation", "University of Florida")

	part, err := writer.CreateFormFile("productImage", "image.jpg")
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}
	_, err = part.Write(file)
	if err != nil {
		t.Fatalf("Error writing file to form: %v", err)
	}

	writer.Close()

	req, _ := http.NewRequest("POST", "/products", &requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	mockProductRepo.On("CreateProduct", mock.AnythingOfType("model.Product")).Return(nil)
	mockImageRepo.On("UploadImage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("test-image-key", nil)

	handler.CreateProductHandler(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockProductRepo.AssertExpectations(t)
	mockImageRepo.AssertExpectations(t)
}

func TestGetAllProductsHandler(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockImageRepo := new(MockImageRepository)
	handler := NewProductHandler(mockProductRepo, mockImageRepo)

	products := []model.Product{
		{UserID: 1, ProductTitle: "Product 1", ProductID: "product1"},
		{UserID: 2, ProductTitle: "Product 2", ProductID: "product2"},
	}

	lastID := "product1"
	limit := 5

	mockProductRepo.On("GetAllProducts", lastID, limit).Return(products, nil)
	mockImageRepo.On("GetPreSignedURLs", mock.Anything).Return(products)

	req, _ := http.NewRequest("GET", "/products?lastID="+lastID+"&limit="+strconv.Itoa(limit), nil)
	rr := httptest.NewRecorder()

	handler.GetAllProductsHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockProductRepo.AssertExpectations(t)
	mockImageRepo.AssertExpectations(t)
}

func TestGetAllProductsByUserIDHandler(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockImageRepo := new(MockImageRepository)
	handler := NewProductHandler(mockProductRepo, mockImageRepo)

	userID := 1
	products := []model.Product{
		{UserID: userID, ProductTitle: "Product 1", ProductID: "product1"},
		{UserID: userID, ProductTitle: "Product 2", ProductID: "product2"},
	}

	lastID := "product1"
	limit := 5

	mockProductRepo.On("GetAllProducts", lastID, limit).Return(products, nil)
	mockImageRepo.On("GetPreSignedURLs", mock.Anything).Return(products)

	req, _ := http.NewRequest("GET", "/products?lastID="+lastID+"&limit="+strconv.Itoa(limit), nil)
	rr := httptest.NewRecorder()

	handler.GetAllProductsHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockProductRepo.AssertExpectations(t)
	mockImageRepo.AssertExpectations(t)
}

func TestUpdateProductHandler(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockImageRepo := new(MockImageRepository)
	handler := NewProductHandler(mockProductRepo, mockImageRepo)

	product := &model.Product{
		UserID:             1,
		ProductID:          "test-product-id",
		ProductTitle:       "Updated Product",
		ProductDescription: "Updated product description",
		ProductPostDate:    "03-03-2025",
		ProductCondition:   4,
		ProductPrice:       199.99,
		ProductLocation:    "New Location",
		ProductImage:       "new-image.png",
	}

	file, err := CreateMockImage("jpeg")
	if err != nil {
		t.Fatalf("Error creating mock image: %v", err)
	}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	_ = writer.WriteField("userId", "1")
	_ = writer.WriteField("productId", "test-product-id")
	_ = writer.WriteField("productTitle", "Updated Product")
	_ = writer.WriteField("productDescription", "Updated product description")
	_ = writer.WriteField("productPostDate", "03-03-2025")
	_ = writer.WriteField("productCondition", "4")
	_ = writer.WriteField("productPrice", "199.99")
	_ = writer.WriteField("productLocation", "New Location")

	part, err := writer.CreateFormFile("productImage", "image.jpg")
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}
	_, err = part.Write(file)
	if err != nil {
		t.Fatalf("Error writing file to form: %v", err)
	}

	writer.Close()

	req, _ := http.NewRequest("PUT", "/products/1/test-product-id", &requestBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	vars := map[string]string{
		"UserId":    "1",
		"ProductId": "test-product-id",
	}
	req = mux.SetURLVars(req, vars)

	mockProductRepo.On("FindProductByUserAndId", 1, "test-product-id").Return(product, nil)
	mockProductRepo.On("UpdateProduct", 1, "test-product-id", mock.AnythingOfType("model.Product")).Return(nil)
	mockImageRepo.On("DeleteImage", mock.Anything).Return(nil)
	mockImageRepo.On("UploadImage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("new-image-key", nil)

	handler.UpdateProductHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockProductRepo.AssertExpectations(t)
	mockImageRepo.AssertExpectations(t)
}

func TestDeleteProductHandler(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockImageRepo := new(MockImageRepository)
	handler := NewProductHandler(mockProductRepo, mockImageRepo)

	userID := 1
	productID := "test-product-id"
	product := &model.Product{
		UserID:       userID,
		ProductID:    productID,
		ProductImage: "test-image-key",
	}

	mockProductRepo.On("FindProductByUserAndId", userID, productID).Return(product, nil)
	mockProductRepo.On("DeleteProduct", userID, productID).Return(nil)
	mockImageRepo.On("DeleteImage", "test-image-key").Return(nil)

	req, _ := http.NewRequest("DELETE", "/products/1/test-product-id", nil)
	rr := httptest.NewRecorder()

	vars := map[string]string{
		"UserId":    "1",
		"ProductId": "test-product-id",
	}
	req = mux.SetURLVars(req, vars)

	handler.DeleteProductHandler(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockProductRepo.AssertExpectations(t)
	mockImageRepo.AssertExpectations(t)
}
