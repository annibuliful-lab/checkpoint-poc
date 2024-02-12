package s3

import (
	"log"
	"os"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	Client     *minio.Client
	BucketName string
}

var s3Singleton *S3Client
var s3Once sync.Once

// GetClient returns a singleton instance of the S3Client.
func GetClient() (*S3Client, error) {
	s3Once.Do(func() {
		accessKey := os.Getenv("S3_ACCESS_KEY")
		secretKey := os.Getenv("S3_SECRET_KEY")
		endpoint := os.Getenv("S3_ENDPOINT")
		bucketName := os.Getenv("S3_BUCKET_NAME")

		minioClient, err := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: false,
		})

		if err != nil {
			log.Fatalln(err)
		}

		s3Singleton = &S3Client{
			Client:     minioClient,
			BucketName: bucketName,
		}
	})

	return s3Singleton, nil
}
