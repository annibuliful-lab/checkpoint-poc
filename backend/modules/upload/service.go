package upload

import (
	"checkpoint/s3"
	"context"
	"time"
)

func signedUrl(ctx context.Context, objectName string) (*string, error) {
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
