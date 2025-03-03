package model

import (
	"errors"
	"testing"
)

func TestProductValidation(t *testing.T) {
	validProduct := Product{
		UserID:             123,
		ProductID:          "9b96a85c-f02e-47a1-9a1a-1dd9ed6147bd",
		ProductTitle:       "Laptop",
		ProductDescription: "A high-performance laptop",
		ProductPostDate:    "02-20-2025",
		ProductCondition:   4,
		ProductPrice:       999.99,
		ProductLocation:    "University of Florida",
		ProductImage:       "https://example.com/laptop.jpg",
	}

	err := validProduct.Validate()
	if err != nil {
		t.Errorf("Expected no error for valid product, but got: %v", err)
	}

	invalidProduct := Product{
		UserID:           0, // Invalid UserID
		ProductID:        "",
		ProductTitle:     "",
		ProductPostDate:  "02-20-2025", // Valid format, but missing other required fields
		ProductCondition: 0,            // Invalid condition
		ProductPrice:     0,            // Invalid price
		ProductLocation:  "",
		ProductImage:     "",
	}

	err = invalidProduct.Validate()
	if err == nil {
		t.Errorf("Expected error for invalid product, but got none")
	}
}

func TestProductPostDateValidation(t *testing.T) {
	productWithValidDate := Product{
		UserID:           123,
		ProductID:        "9b96a85c-f02e-47a1-9a1a-1dd9ed6147bd",
		ProductTitle:     "Laptop",
		ProductPostDate:  "02-20-2025", // Valid date format
		ProductCondition: 4,
		ProductPrice:     999.99,
		ProductLocation:  "University of Florida",
		ProductImage:     "https://example.com/laptop.jpg",
	}

	err := productWithValidDate.Validate()
	if err != nil {
		t.Errorf("Expected no error for valid product post date, but got: %v", err)
	}

	productWithInvalidDate := Product{
		UserID:           123,
		ProductID:        "9b96a85c-f02e-47a1-9a1a-1dd9ed6147bd",
		ProductTitle:     "Laptop",
		ProductPostDate:  "2025-02-20", // Invalid date format
		ProductCondition: 4,
		ProductPrice:     999.99,
		ProductLocation:  "University of Florida",
		ProductImage:     "https://example.com/laptop.jpg",
	}

	err = productWithInvalidDate.Validate()
	if err == nil || err.Error() != "validation failed: productPostDate must be in MM-DD-YYYY format" {
		t.Errorf("Expected error for invalid product post date, but got: %v", err)
	}
}

func TestProductValidationWithEmptyFields(t *testing.T) {
	invalidProduct := Product{
		UserID:           0, // Invalid UserID
		ProductID:        "",
		ProductTitle:     "",           // Invalid ProductTitle
		ProductPostDate:  "02-20-2025", // Valid date format, but missing other required fields
		ProductCondition: 0,            // Invalid ProductCondition
		ProductPrice:     0,            // Invalid ProductPrice
		ProductLocation:  "",
		ProductImage:     "",
	}

	err := invalidProduct.Validate()
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
