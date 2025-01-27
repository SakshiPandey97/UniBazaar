package model

import (
	"fmt"

	"gopkg.in/validator.v2"
)

type Product struct {
	ProductID   string  `json:"productId" bson:"productId"`
	Name        string  `json:"name" bson:"name" validate:"nonzero"`
	Price       float64 `json:"price" bson:"price" validate:"nonzero"`
	Age         int     `json:"age" bson:"age" validate:"nonzero"`
	Description string  `json:"description" bson:"description"`
}

func (p *Product) Validate() error {
	if err := validator.Validate(p); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}
	return nil
}
