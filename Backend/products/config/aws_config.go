package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var awsClientObj *s3.Client

type Loader interface {
	LoadDefaultConfig(ctx context.Context, optFns ...func(*awsConfig.LoadOptions) error) (aws.Config, error)
}

type DefaultLoader struct{}

func (l DefaultLoader) LoadDefaultConfig(ctx context.Context, optFns ...func(*awsConfig.LoadOptions) error) (aws.Config, error) {
	return awsConfig.LoadDefaultConfig(ctx, optFns...)
}

func loadAWSConfig(ctx context.Context, loader Loader) (*s3.Client, error) {
	cfg, err := loader.LoadDefaultConfig(ctx)
	// Uncomment to debug
	// config.WithClientLogMode(aws.LogRequest|aws.LogResponse|aws.LogRetries),
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	return s3.NewFromConfig(cfg), nil
}

func GetAWSClientInstance(ctx context.Context, loader Loader) (*s3.Client, error) {
	if awsClientObj != nil {
		return awsClientObj, nil
	}

	client, err := loadAWSConfig(ctx, loader)
	if err != nil {
		return nil, err
	}

	awsClientObj = client
	return awsClientObj, nil
}
