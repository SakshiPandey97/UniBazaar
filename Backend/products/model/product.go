package model

import (
	"fmt"
	"regexp"

	"gopkg.in/validator.v2"
)

type UserProduct struct {
	UserID   string    `json:"UserId" bson:"UserId"`
	Products []Product `json:"Products" bson:"Products"`
}

type Product struct {
	ProductID          string  `json:"ProductId" bson:"ProductId"`
	ProductTitle       string  `json:"ProductTitle" bson:"ProductTitle" validate:"nonzero"`
	ProductDescription string  `json:"ProductDescription" bson:"ProductDescription"`
	ProductPostDate    string  `json:"ProductPostDate" bson:"ProductPostDate" validate:"nonzero"`
	ProductCondition   int     `json:"ProductCondition" bson:"ProductCondition" validate:"nonzero"`
	ProductPrice       float64 `json:"ProductPrice" bson:"ProductPrice" validate:"nonzero"`
	ProductLocation    string  `json:"ProductLocation" bson:"ProductLocation"`
	ProductImage       string  `json:"ProductImage" bson:"ProductImage" validate:"required"`
}

func (p *Product) Validate() error {
	if err := validator.Validate(p); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

	dateRegex := `^\d{2}-\d{2}-\d{4}$`
	matched, err := regexp.MatchString(dateRegex, p.ProductPostDate)
	if err != nil {
		return fmt.Errorf("error while validating date: %v", err)
	}
	if !matched {
		return fmt.Errorf("validation failed: productPostDate must be in YYYY-MM-DD format")
	}

	return nil
}
