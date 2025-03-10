package model

import (
	"errors"
	"testing"
	"time"
)

func TestProductValidation(t *testing.T) {
	productPostDate, err := time.Parse("01-02-2006", "03-03-2025")
	if err != nil {
		t.Fatalf("Failed to parse date: %v", err)
	}

	validProduct := Product{
		UserID:             123,
		ProductID:          "9b96a85c-f02e-47a1-9a1a-1dd9ed6147bd",
		ProductTitle:       "Laptop",
		ProductDescription: "A high-performance laptop",
		ProductPostDate:    productPostDate,
		ProductCondition:   4,
		ProductPrice:       999.99,
		ProductLocation:    "University of Florida",
		ProductImage:       "https://example.com/laptop.jpg",
	}

	err = validProduct.Validate()
	if err != nil {
		t.Errorf("Expected no error for valid product, but got: %v", err)
	}

	invalidProduct := Product{
		UserID:           0, // Invalid UserID
		ProductID:        "",
		ProductTitle:     "",
		ProductPostDate:  productPostDate,
		ProductCondition: 0, // Invalid condition
		ProductPrice:     0, // Invalid price
		ProductLocation:  "",
		ProductImage:     "",
	}

	err = invalidProduct.Validate()
	if err == nil {
		t.Errorf("Expected error for invalid product, but got none")
	}
}

func TestProductValidationWithEmptyFields(t *testing.T) {
	productPostDate, err := time.Parse("01-02-2006", "03-03-2025")
	if err != nil {
		t.Fatalf("Failed to parse date: %v", err)
	}
	invalidProduct := Product{
		UserID:           0, // Invalid UserID
		ProductID:        "",
		ProductTitle:     "", // Invalid ProductTitle
		ProductPostDate:  productPostDate,
		ProductCondition: 0, // Invalid ProductCondition
		ProductPrice:     0, // Invalid ProductPrice
		ProductLocation:  "",
		ProductImage:     "",
	}

	err = invalidProduct.Validate()
	if err == nil {
		t.Errorf("Expected error for product with empty required fields, but got none")
	}
}

func TestFormatValidationError(t *testing.T) {
	originalError := errors.New("ProductTitle: zero value")
	formattedError := formatValidationError(originalError)

	if formattedError.Error() != "ProductTitle: cannot be empty or zero" {
		t.Errorf("Expected 'ProductTitle: cannot be empty or zero', but got: %v", formattedError)
	}
}

func TestEmptyProduct(t *testing.T) {
	emptyProduct := Product{}

	err := emptyProduct.Validate()
	if err == nil {
		t.Errorf("Expected error for empty product, but got none")
	}
}

func TestFormatValidationErrorWithNilError(t *testing.T) {
	err := formatValidationError(nil)
	if err != nil {
		t.Errorf("Expected nil error, but got: %v", err)
	}
}
