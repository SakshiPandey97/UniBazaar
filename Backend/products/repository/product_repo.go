package repository

import "web-service/model"

type ProductRepository interface {
	CreateProduct(product model.Product) error
	GetAllProducts() ([]model.Product, error)
	GetProductsByUserID(userID int) ([]model.Product, error)
	UpdateProduct(userID int, productID string, product model.Product) error
	DeleteProduct(userID int, productID string) error
	FindProductByUserAndId(userID int, productID string) (*model.Product, error)
}
