package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	args := m.Called(append([]interface{}{ctx, params}, convertToInterfaceSlice(optFns)...)...)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*s3.ListBucketsOutput), args.Error(1)
}

func convertToInterfaceSlice(fns []func(*s3.Options)) []interface{} {
	result := make([]interface{}, len(fns))
	for i, fn := range fns {
		result[i] = fn
	}
	return result
}

type MockConfigLoader struct {
	mock.Mock
}

func (m *MockConfigLoader) LoadDefaultConfig(ctx context.Context, optFns ...func(*awsConfig.LoadOptions) error) (aws.Config, error) {
	interfaceOptFns := make([]interface{}, len(optFns))
	for i, fn := range optFns {
		interfaceOptFns[i] = fn
	}

	args := m.Called(append([]interface{}{ctx}, interfaceOptFns...)...)

	if args.Get(0) == nil {
		fmt.Println("MockConfigLoader: returning error")
		return aws.Config{}, args.Error(1)
	}
	fmt.Println("MockConfigLoader: returning success")
	return args.Get(0).(aws.Config), args.Error(1)
}

func TestGetAWSClientInstance_Failure(t *testing.T) {
	mockLoader := new(MockConfigLoader)
	mockLoader.On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{}, assert.AnError)

	client, err := GetAWSClientInstance(context.Background(), mockLoader)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestGetAWSClientInstance_Success(t *testing.T) {
	mockLoader := new(MockConfigLoader)
	cfg := aws.Config{Region: "us-east-1"}
	mockLoader.On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(cfg, nil)

	client, err := GetAWSClientInstance(context.Background(), mockLoader)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "us-east-1", client.Options().Region)
}

func TestGetAWSClientInstance_Singleton(t *testing.T) {
	mockLoader := new(MockConfigLoader)
	cfg := aws.Config{Region: "us-east-1"}
	mockLoader.On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(cfg, nil)

	client1, err1 := GetAWSClientInstance(context.Background(), mockLoader)
	assert.NoError(t, err1)
	assert.NotNil(t, client1)

	client2, err2 := GetAWSClientInstance(context.Background(), mockLoader)
	assert.NoError(t, err2)
	assert.NotNil(t, client2)

	assert.Same(t, client1, client2)
}

func TestLoadAWSConfig_Error(t *testing.T) {
	mockLoader := new(MockConfigLoader)
	mockLoader.On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(aws.Config{}, assert.AnError)

	client, err := loadAWSConfig(context.Background(), mockLoader)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestLoadAWSConfig_Success(t *testing.T) {
	mockLoader := new(MockConfigLoader)
	cfg := aws.Config{Region: "us-east-1"}
	mockLoader.On("LoadDefaultConfig", mock.Anything, mock.Anything).Return(cfg, nil)

	client, err := loadAWSConfig(context.Background(), mockLoader)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "us-east-1", client.Options().Region)
}

func TestConvertToInterfaceSlice_Empty(t *testing.T) {
	result := convertToInterfaceSlice([]func(*s3.Options){})
	assert.Empty(t, result)
}

func TestConvertToInterfaceSlice_NonEmpty(t *testing.T) {
	fn := func(options *s3.Options) {
		options.UsePathStyle = true
	}

	result := convertToInterfaceSlice([]func(*s3.Options){fn})
	assert.Len(t, result, 1)
	assert.NotNil(t, result[0])
}
