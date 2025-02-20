package repository

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
	"web-service/config"
	"web-service/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadToS3Bucket(productId string, userId string, fileData []byte, filetype string) (string, error) {
	awsClientObj := config.GetAWSClientInstance()
	uploader := manager.NewUploader(awsClientObj)

	fileParts := strings.Split(filetype, "/")
	if len(fileParts) != 2 {
		return "", fmt.Errorf("invalid filetype format: %s", filetype)
	}
	fileExtension := fileParts[1]

	objectKey := fmt.Sprintf("products/%s/%s.%s", userId, productId, fileExtension)

	_, uploadError := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("unibazaar-bucket"),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(fileData),
	})
	if uploadError != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", uploadError)
	}

	return objectKey, nil
}

func DeleteImageFromS3(objectKey string) error {
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

func generatePresignedURL(objectKey string) (string, error) {
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

func GetPreSignedURLs(products []model.Product) []model.Product {
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := range products {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			preSignedURL, err := generatePresignedURL(products[i].ProductImage)
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
