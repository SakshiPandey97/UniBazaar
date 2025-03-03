package model

import (
	"encoding/json"
	"testing"
)

func TestErrorResponseSerialization(t *testing.T) {
	errorResponse := ErrorResponse{
		Error:   "Error updating product",
		Details: "ProductPrice: cannot be empty or zero, Product not found",
	}

	jsonData, err := json.Marshal(errorResponse)
	if err != nil {
		t.Fatalf("Error marshalling ErrorResponse: %v", err)
	}

	expectedJSON := `{"error":"Error updating product","details":"ProductPrice: cannot be empty or zero, Product not found"}`
	if string(jsonData) != expectedJSON {
		t.Errorf("Expected JSON: %s, but got: %s", expectedJSON, string(jsonData))
	}
}

func TestErrorResponseDeserialization(t *testing.T) {
	jsonData := `{"error":"Error updating product","details":"ProductPrice: cannot be empty or zero, Product not found"}`

	var errorResponse ErrorResponse

	err := json.Unmarshal([]byte(jsonData), &errorResponse)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON to ErrorResponse: %v", err)
	}

	if errorResponse.Error != "Error updating product" {
		t.Errorf("Expected error message: 'Error updating product', but got: %s", errorResponse.Error)
	}
	if errorResponse.Details != "ProductPrice: cannot be empty or zero, Product not found" {
		t.Errorf("Expected details: 'ProductPrice: cannot be empty or zero, Product not found', but got: %s", errorResponse.Details)
	}
}

func TestErrorResponseSerializationWithoutDetails(t *testing.T) {
	errorResponse := ErrorResponse{
		Error: "Error updating product",
	}

	jsonData, err := json.Marshal(errorResponse)
	if err != nil {
		t.Fatalf("Error marshalling ErrorResponse: %v", err)
	}

	expectedJSON := `{"error":"Error updating product"}`
	if string(jsonData) != expectedJSON {
		t.Errorf("Expected JSON: %s, but got: %s", expectedJSON, string(jsonData))
	}
}

func TestEmptyErrorResponse(t *testing.T) {
	errorResponse := ErrorResponse{}

	jsonData, err := json.Marshal(errorResponse)
	if err != nil {
		t.Fatalf("Error marshalling ErrorResponse: %v", err)
	}

	expectedJSON := `{"error":""}`
	if string(jsonData) != expectedJSON {
		t.Errorf("Expected JSON: %s, but got: %s", expectedJSON, string(jsonData))
	}
}
