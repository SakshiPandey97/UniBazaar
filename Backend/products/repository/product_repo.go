package repository

import "web-service/model"

type ProductRepository interface {
	CreateProduct(product model.Product) error
	GetAllProducts(lastId string, limit int) ([]model.Product, error)
	GetProductsByUserID(userID int, lastId string, limit int) ([]model.Product, error)
	UpdateProduct(userID int, productID string, product model.Product) error
	DeleteProduct(userID int, productID string) error
	FindProductByUserAndId(userID int, productID string) (*model.Product, error)
	SearchProducts(query string, limit int) ([]model.Product, error)
}
