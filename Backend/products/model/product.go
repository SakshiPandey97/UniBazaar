package model

import (
	"errors"
	"strings"
	"time"

	"gopkg.in/validator.v2"
)

// Product represents a product in the marketplace.
// @Description Represents a product for sale in the marketplace.
// @Type Product
// @Property userId int "Unique user ID" required example(123)
// @Property productId string "Unique product ID" required example("9b96a85c-f02e-47a1-9a1a-1dd9ed6147bd")
// @Property productTitle string "Product title" required example("Laptop")
// @Property productDescription string "Product description" example("A high-performance laptop")
// @Property productPostDate string "Product post date in MM-DD-YYYY format" required format("date") example("02-20-2025")
// @Property productCondition int "Product condition" required example(4)
// @Property productPrice float64 "Price of the product" required example(999.99)
// @Property productLocation string "Location of the product" example("University of Florida")
// @Property productImage string "In POST: The product image file. In GET: The URL of the product image" example("https://example.com/laptop.jpg")
type Product struct {
	UserID             int       `json:"userId" bson:"UserId" validate:"nonzero" example:"123"`                            // Unique user ID
	ProductID          string    `json:"productId" bson:"ProductId" example:"9b96a85c-f02e-47a1-9a1a-1dd9ed6147bd"`        // Unique product ID (UUID)
	ProductTitle       string    `json:"productTitle" bson:"ProductTitle" validate:"nonzero" example:"Laptop"`             // Product title
	ProductDescription string    `json:"productDescription" bson:"ProductDescription" example:"A high-performance laptop"` // Product description
	ProductPostDate    time.Time `json:"productPostDate" bson:"ProductPostDate" validate:"nonzero" example:"02-20-2025"`   // Product post date (time.Time)
	ProductCondition   int       `json:"productCondition" bson:"ProductCondition" validate:"nonzero" example:"4"`          // Product condition
	ProductPrice       float64   `json:"productPrice" bson:"ProductPrice" validate:"nonzero" example:"999.99"`             // Price of the product
	ProductLocation    string    `json:"productLocation" bson:"ProductLocation" example:"University of Florida"`           // Location of the product
	ProductImage       string    `json:"productImage" bson:"ProductImage" example:"https://example.com/laptop.jpg"`        // Product image URL in GET, Actual product image in PUT
}

func (p *Product) Validate() error {
	if err := validator.Validate(p); err != nil {
		return formatValidationError(err)
	}

	return nil
}

func formatValidationError(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	replacements := map[string]string{
		"zero value": "cannot be empty or zero",
	}

	for old, new := range replacements {
		errMsg = strings.ReplaceAll(errMsg, old, new)
	}

	return errors.New(errMsg)
}
