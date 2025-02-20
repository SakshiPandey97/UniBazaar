package config

import (
	"context"
	"log"

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

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // Enforce path-style addressing
	})
}

func GetAWSClientInstance() *s3.Client {
	if awsClientObj != nil {
		return awsClientObj
	} else {
		awsClientObj = loadAWSConfig()
		return awsClientObj
	}
}
