package helper

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"web-service/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadToS3Bucket(productId string, userId string, fileData []byte, filetype string) (string, error) {
	var awsClientObj = config.GetAWSClientInstance()
	uploader := manager.NewUploader(awsClientObj)

	objectKey := fmt.Sprintf("products/%s/%s.%s", userId, productId, strings.Split(filetype, "/")[1])

	_, uploadError := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("unibazaar-bucket"),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(fileData),
	})
	if uploadError != nil {
		log.Printf("Error uploading to Amazon S3: %v\n", uploadError)
		return "", uploadError
	}
	return objectKey, nil
}
