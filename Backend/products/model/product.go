package model

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/validator.v2"
)

type UserProduct struct {
	UserID   int       `json:"userId" bson:"UserId" validate:"nonzero"`
	Products []Product `json:"products" bson:"Products"`
}

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
