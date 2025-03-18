package config

import (
	"context"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*s3.ListBucketsOutput), args.Error(1)
}

type MockConfigLoader struct {
	mock.Mock
}

func (m *MockConfigLoader) LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	args := m.Called(ctx, optFns)
	return args.Get(0).(aws.Config), args.Error(1)
}

var loadConfigFunc = func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, optFns...)
}

func TestLoadAWSConfig_Success(t *testing.T) {
	client := loadAWSConfig()
	if client == nil {
		t.Error("loadAWSConfig returned nil")
	}

	mockS3 := new(MockS3Client)
	mockS3.On("ListBuckets", mock.Anything, mock.Anything, mock.Anything).Return(&s3.ListBucketsOutput{
		Buckets: []types.Bucket{
			{Name: aws.String("test-bucket")},
		},
		ResultMetadata: middleware.Metadata{},
	}, nil)

	if client != nil {
		_, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
		if err != nil && awsClientObj == nil {
			log.Println("Error listing buckets: ", err)
		}
	}
}

func TestGetAWSClientInstance_SingleInstance(t *testing.T) {
	client1 := GetAWSClientInstance()
	client2 := GetAWSClientInstance()

	if client1 != client2 {
		t.Error("GetAWSClientInstance returned different instances")
	}

	if client1 == nil {
		t.Error("GetAWSClientInstance returned nil")
	}
}

func TestGetAWSClientInstance_MultipleCalls(t *testing.T) {
	client1 := GetAWSClientInstance()
	client2 := GetAWSClientInstance()

	assert.Same(t, client1, client2, "should return the same instance")
}
