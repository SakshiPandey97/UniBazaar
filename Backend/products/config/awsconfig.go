package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var awsClientObj *s3.Client

func loadAWSConfig() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return nil
	}

	return s3.NewFromConfig(cfg)
}

func GetAWSClientInstance() *s3.Client {
	if awsClientObj != nil {
		return awsClientObj
	} else {
		awsClientObj = loadAWSConfig()
		return awsClientObj
	}
}

func GeneratePresignedURL(objectKey string) (string, error) {
	client := GetAWSClientInstance()
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
