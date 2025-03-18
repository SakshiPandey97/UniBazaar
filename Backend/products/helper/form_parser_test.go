package helper

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/url"
	"testing"
	customerrors "web-service/errors"
	"web-service/model"
)

func TestGetUserID_ValidInput(t *testing.T) {
	result, err := GetUserID("123")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 123 {
		t.Errorf("Expected: 123, got: %d", result)
	}
}

func TestGetUserID_InvalidInput(t *testing.T) {
	tests := []string{"abc", "", "12.34", "-"}

	for _, input := range tests {
		t.Run("Invalid Input: "+input, func(t *testing.T) {
			_, err := GetUserID(input)
			if err == nil {
				t.Errorf("Expected error for input '%s', but got nil", input)
			}
			if _, ok := err.(*customerrors.BadRequestError); !ok {
				t.Errorf("Expected BadRequestError for input '%s', but got %T", input, err)
			}
		})
	}
}

func TestParseFormAndCreateProduct_ValidData(t *testing.T) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	_ = writer.WriteField("productTitle", "Sample Product")
	_ = writer.WriteField("productDescription", "A great product")
	_ = writer.WriteField("productPostDate", "03-03-2025")
	_ = writer.WriteField("productLocation", "New York")
	_ = writer.WriteField("productImage", "image.png")
	_ = writer.WriteField("productCondition", "1")
	_ = writer.WriteField("productPrice", "99.99")

	writer.Close()

	req, err := http.NewRequest("POST", "/", &requestBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	userID := 1001
	product, err := ParseFormAndCreateProduct(req, userID)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if product.UserID != userID {
		t.Errorf("Expected UserID: %d, got: %d", userID, product.UserID)
	}

	if product.ProductTitle != "Sample Product" {
		t.Errorf("Expected title: 'Sample Product', got: '%s'", product.ProductTitle)
	}

	if product.ProductID == "" {
		t.Errorf("ProductID should not be empty, got: %v", product.ProductID)
	}

	if product.ProductPostDate.Format("01-02-2006") != "03-03-2025" {
		t.Errorf("Expected post date: 03-03-2025, got: %v", product.ProductPostDate.Format("01-02-2006"))
	}

	if product.ProductCondition != 1 {
		t.Errorf("Expected condition: 1, got: %d", product.ProductCondition)
	}

	if product.ProductPrice != 99.99 {
		t.Errorf("Expected price: 99.99, got: %f", product.ProductPrice)
	}
}

func TestParseFormAndCreateProduct_MissingOrInvalidData(t *testing.T) {
	tests := []struct {
		name     string
		formData url.Values
	}{
		{"Missing Product Title", url.Values{"productDescription": {"A great product"}, "productPostDate": {"03-03-2025"}}},
		{"Invalid Product Condition", url.Values{"productCondition": {"not_a_number"}, "productPostDate": {"03-03-2025"}}},
		{"Invalid Product Price", url.Values{"productPrice": {"not_a_number"}, "productPostDate": {"03-03-2025"}}},
		{"Missing Product Post Date", url.Values{"productTitle": {"title"}}},
		{"Invalid Product Post Date", url.Values{"productPostDate": {"33-33-3333"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{Form: tt.formData}
			_, err := ParseFormAndCreateProduct(req, 1001)

			if err == nil {
				t.Errorf("Expected error for case '%s', but got nil", tt.name)
			}
			if _, ok := err.(*customerrors.BadRequestError); !ok {
				t.Errorf("Expected BadRequestError for input '%s', but got %T", tt.name, err)
			}
		})
	}
}

func TestParseNumericalFormValues_ValidData(t *testing.T) {
	formData := url.Values{
		"productCondition": {"2"},
		"productPrice":     {"49.99"},
	}

	req := &http.Request{Form: formData}
	product := model.Product{}

	err := parseNumericalFormValues(req, &product)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if product.ProductCondition != 2 {
		t.Errorf("Expected condition: 2, got: %d", product.ProductCondition)
	}

	if product.ProductPrice != 49.99 {
		t.Errorf("Expected price: 49.99, got: %f", product.ProductPrice)
	}
}

func TestParseNumericalFormValues_InvalidData(t *testing.T) {
	tests := []struct {
		name     string
		formData url.Values
	}{
		{"Invalid Product Condition", url.Values{"productCondition": {"invalid"}}},
		{"Invalid Product Price", url.Values{"productPrice": {"invalid"}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{Form: tt.formData}
			product := model.Product{}

			err := parseNumericalFormValues(req, &product)
			if err == nil {
				t.Errorf("Expected error for case '%s', but got nil", tt.name)
			}
		})
	}
}

func TestParseLimit(t *testing.T) {
	tests := []struct {
		name     string
		limitStr string
		expected int
	}{
		{"Valid limit", "20", 20},
		{"Empty limit", "", 10},
		{"Invalid limit (non-numeric)", "abc", 10},
		{"Negative limit", "-5", 10},
		{"Zero limit", "0", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseLimit(tt.limitStr)
			if result != tt.expected {
				t.Errorf("For input '%s', expected %d, but got %d", tt.limitStr, tt.expected, result)
			}
		})
	}
}

func TestCheckParam_Valid(t *testing.T) {
	param := "validParam"
	result, err := CheckParam(param)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != param {
		t.Errorf("Expected '%s', got '%s'", param, result)
	}
}

func TestCheckParam_Empty(t *testing.T) {
	_, err := CheckParam("")
	if err == nil {
		t.Errorf("Expected error for empty param, but got nil")
	}
	if _, ok := err.(*customerrors.BadRequestError); !ok {
		t.Errorf("Expected BadRequestError, got %T", err)
	}
}
