package repository

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"web-service/config"
	"web-service/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3ImageRepository struct{}

func NewS3ImageRepository() *S3ImageRepository {
	return &S3ImageRepository{}
}

func (r *S3ImageRepository) UploadImage(productId string, userId string, fileData []byte, filetype string) (string, error) {
	start := time.Now()

	awsClientObj := config.GetAWSClientInstance()
	uploader := manager.NewUploader(awsClientObj)

	objectKey := fmt.Sprintf("products/%s/%s.%s", userId, productId, filetype)
	log.Println("Uploading to S3 with key:", objectKey)

	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("unibazaar-bucket"),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(fileData),
	})
	if err != nil {
		log.Printf("Error uploading to Amazon S3: %v\n", err)
		return "", err
	}

	duration := time.Since(start)
	log.Printf("UploadToS3Bucket took %s\n", duration)

	return objectKey, nil
}

func (r *S3ImageRepository) DeleteImage(objectKey string) error {
	client := config.GetAWSClientInstance()

	_, err := client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String("unibazaar-bucket"),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}

	log.Printf("Successfully deleted object %s from S3", objectKey)
	return nil
}

func (r *S3ImageRepository) GeneratePresignedURL(objectKey string) (string, error) {
	client := config.GetAWSClientInstance()
	log.Printf("Pre Key %s", objectKey)
	psClient := s3.NewPresignClient(client)
	req, err := psClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("unibazaar-bucket"),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(15 * int64(time.Minute))
	}) // URL expires in 15 min

	if err != nil {
		return "", fmt.Errorf("failed to generate pre-signed URL: %w", err)
	}
	return req.URL, nil
}

func (r *S3ImageRepository) GetPreSignedURLs(products []model.Product) []model.Product {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := range products {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			preSignedURL, err := r.GeneratePresignedURL(products[i].ProductImage)
			if err != nil {
				log.Printf("Failed to generate pre-signed URL for ProductID %s: %v", products[i].ProductID, err)
				return
			}

			mu.Lock()
			products[i].ProductImage = preSignedURL
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	return products
}
