package repository

import (
	"context"
	"errors"
	"testing"

	"web-service/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
	args := m.Called(ctx, input)
	return &manager.UploadOutput{}, args.Error(1)
}

func (m *MockS3Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	args := m.Called(ctx, input)
	return &s3.DeleteObjectOutput{}, args.Error(1)
}

func (m *MockS3Client) PresignGetObject(ctx context.Context, input *s3.GetObjectInput, opts ...func(*s3.PresignOptions)) (string, error) {
	args := m.Called(ctx, input)
	return "https://mock-url.com", args.Error(1)
}

type MockAWSClientFactory struct {
	mock.Mock
}

func (m *MockAWSClientFactory) GetAWSClientInstance() aws.Config {
	args := m.Called()
	return args.Get(0).(aws.Config)
}

// Test UploadImage Success
// func TestUploadImage_Success(t *testing.T) {
// 	mockS3 := new(MockS3Client)
// 	mockFactory := new(MockAWSClientFactory)

// 	// Setup the mock to return a successful response for S3 upload
// 	mockS3.On("Upload", mock.Anything, mock.Anything).Return(&manager.UploadOutput{}, nil)
// 	mockFactory.On("GetAWSClientInstance").Return(aws.Config{})

// 	repo := &S3ImageRepository{}

// 	// Call the method under test
// 	fileData := []byte("mock image data")
// 	url, err := repo.UploadImage("123", "456", fileData, "jpg")

// 	// Verify the results
// 	assert.NoError(t, err)
// 	assert.Contains(t, url, "products/456/123.jpg")
// }

func TestUploadImage_Failure(t *testing.T) {
	mockS3 := new(MockS3Client)
	mockFactory := new(MockAWSClientFactory)

	mockS3.On("Upload", mock.Anything, mock.Anything).Return(&manager.UploadOutput{}, errors.New("upload error"))
	mockFactory.On("GetAWSClientInstance").Return(aws.Config{})

	repo := &S3ImageRepository{}

	fileData := []byte("mock image data")
	url, err := repo.UploadImage("123", "456", fileData, "jpg")

	assert.Error(t, err)
	assert.Empty(t, url)
}

// Test DeleteImage Success
// func TestDeleteImage_Success(t *testing.T) {
// 	mockS3 := new(MockS3Client)
// 	mockFactory := new(MockAWSClientFactory)

// 	// Setup the mock to return a successful response for S3 delete
// 	mockS3.On("DeleteObject", mock.Anything, mock.Anything).Return(&s3.DeleteObjectOutput{}, nil)
// 	mockFactory.On("GetAWSClientInstance").Return(aws.Config{})

// 	repo := &S3ImageRepository{}

// 	// Call the method under test
// 	err := repo.DeleteImage("products/456/123.jpg")

// 	// Verify the results
// 	assert.NoError(t, err)
// }

func TestDeleteImage_Failure(t *testing.T) {
	mockS3 := new(MockS3Client)
	mockFactory := new(MockAWSClientFactory)

	mockS3.On("DeleteObject", mock.Anything, mock.Anything).Return(&s3.DeleteObjectOutput{}, errors.New("delete error"))
	mockFactory.On("GetAWSClientInstance").Return(aws.Config{})

	repo := &S3ImageRepository{}

	err := repo.DeleteImage("products/456/123.jpg")

	assert.Error(t, err)
}

// Test GeneratePresignedURL Success
// func TestGeneratePresignedURL_Success(t *testing.T) {
// 	mockS3 := new(MockS3Client)
// 	mockFactory := new(MockAWSClientFactory)

// 	// Setup the mock to return a presigned URL string
// 	mockS3.On("PresignGetObject", mock.Anything, mock.Anything).Return("https://mock-url.com", nil)
// 	mockFactory.On("GetAWSClientInstance").Return(aws.Config{})

// 	repo := &S3ImageRepository{}

// 	// Call the method under test
// 	url, err := repo.GeneratePresignedURL("products/456/123.jpg")

// 	// Verify the results
// 	assert.NoError(t, err)
// 	assert.Equal(t, "https://mock-url.com", url)
// }

func TestGeneratePresignedURL_Failure(t *testing.T) {
	mockS3 := new(MockS3Client)
	mockFactory := new(MockAWSClientFactory)

	mockS3.On("PresignGetObject", mock.Anything, mock.Anything).Return("", errors.New("presign error"))
	mockFactory.On("GetAWSClientInstance").Return(aws.Config{})

	repo := &S3ImageRepository{}

	url, err := repo.GeneratePresignedURL("products/456/123.jpg")

	assert.Error(t, err)
	assert.Empty(t, url)
}

func TestGetPreSignedURLs_Success(t *testing.T) {
	mockS3 := new(MockS3Client)
	mockFactory := new(MockAWSClientFactory)

	mockS3.On("PresignGetObject", mock.Anything, mock.Anything).Return("https://mock-url.com", nil)
	mockFactory.On("GetAWSClientInstance").Return(aws.Config{})

	products := []model.Product{
		{UserID: 456, ProductImage: "product1.jpg"},
		{UserID: 789, ProductImage: "product2.jpg"},
	}

	repo := &S3ImageRepository{}

	result := repo.GetPreSignedURLs(products)

	assert.Len(t, result, 2)

}
