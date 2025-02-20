package model

// ErrorResponse represents an error response.
// @Description Represents an error response when an operation fails.
// @Type ErrorResponse
// @Property error string "Error message"
// @Property details string "Detailed error message"
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
