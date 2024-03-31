package upload

import (
	"checkpoint/s3"
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

func UploadFile(ctx context.Context, objectName string, file io.Reader, size int64) (*string, error) {
	minioClient, err := s3.GetClient()

	if err != nil {
		return nil, err
	}

	minioClient.Client.PutObject(ctx, minioClient.BucketName, objectName, file, size, minio.PutObjectOptions{})
	url, err := SignedUrl(ctx, objectName)

	return url, err
}

func SignedUrl(ctx context.Context, objectName string) (*string, error) {
	minioClient, err := s3.GetClient()

	if err != nil {
		return nil, err
	}

	// Generate a signed URL for the object
	presignedURL, err := minioClient.Client.PresignedGetObject(ctx, minioClient.BucketName, objectName, 15*time.Minute, nil)

	if err != nil {
		return nil, err
	}
	url := presignedURL.String()

	return &url, err
}
