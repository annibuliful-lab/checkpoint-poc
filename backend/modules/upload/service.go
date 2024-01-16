package upload

import (
	"checkpoint/s3"
	"context"
	"time"
)

func signedUrl(ctx context.Context, objectName string) (string, error) {
	minioClient, bucketName, err := s3.GetClient()

	if err != nil {
		return "", err
	}

	// Generate a signed URL for the object
	presignedURL, err := minioClient.PresignedPutObject(ctx, bucketName, objectName, 15*time.Minute)

	if err != nil {
		return "", err
	}

	return presignedURL.String(), err
}
