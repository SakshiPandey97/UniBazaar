package model

// ErrorResponse represents an error response.
// @Description Represents an error response when an operation fails.
// @Type ErrorResponse
// @Property error string "Error message" example("Error updaing product")
// @Property details string "Detailed error message" example("ProductPrice: cannot be empty or zero, Product not found")
type ErrorResponse struct {
	Error   string `json:"error" example:"Error updaing product"`                                                // Error message with example value
	Details string `json:"details,omitempty" example:"ProductPrice: cannot be empty or zero, Product not found"` // Detailed error message with multiple examples
}
