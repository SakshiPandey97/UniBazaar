package repository

import "web-service/model"

type ImageRepository interface {
	UploadImage(productId string, userId string, fileData []byte, filetype string) (string, error)
	DeleteImage(objectKey string) error
	GeneratePresignedURL(objectKey string) (string, error)
	GetPreSignedURLs(products []model.Product) []model.Product
}
