package model

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/validator.v2"
)

// UserProduct represents a user's product list.
// @Description Represents a user along with their associated products in the marketplace.
// @Type UserProduct
// @Property userId int "Unique user ID" required
// @Property products array "List of products owned by the user" required
type UserProduct struct {
	UserID   int       `json:"userId" bson:"UserId" validate:"nonzero"`
	Products []Product `json:"products" bson:"Products"`
}

// Product represents a product in the marketplace.
// @Description Represents a product for sale in the marketplace.
// @Type Product
// @Property productId string "Unique product ID" required
// @Property productTitle string "Product title" required
// @Property productDescription string "Product description"
// @Property productPostDate string "Product post date in DD-MM-YYYY format" required
// @Property productCondition int "Product condition" required
// @Property productPrice float64 "Price of the product" required
// @Property productLocation string "Location of the product"
// @Property productImage string "URL or key of the product image"
type Product struct {
	ProductID          string  `json:"productId" bson:"ProductId"`
	ProductTitle       string  `json:"productTitle" bson:"ProductTitle" validate:"nonzero"`
	ProductDescription string  `json:"productDescription" bson:"ProductDescription"`
	ProductPostDate    string  `json:"productPostDate" bson:"ProductPostDate" validate:"nonzero"`
	ProductCondition   int     `json:"productCondition" bson:"ProductCondition" validate:"nonzero"`
	ProductPrice       float64 `json:"productPrice" bson:"ProductPrice" validate:"nonzero"`
	ProductLocation    string  `json:"productLocation" bson:"ProductLocation"`
	ProductImage       string  `json:"productImage" bson:"ProductImage"`
}

func (p *Product) Validate() error {
	if err := validator.Validate(p); err != nil {
		return formatValidationError(err)
	}

	dateRegex := `^\d{2}-\d{2}-\d{4}$`
	matched, err := regexp.MatchString(dateRegex, p.ProductPostDate)
	if err != nil {
		return fmt.Errorf("error while validating date: %v", err)
	}
	if !matched {
		return fmt.Errorf("validation failed: productPostDate must be in DD-MM-YYYY format")
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
