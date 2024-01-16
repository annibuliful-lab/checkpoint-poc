package s3

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func GetClient() (*minio.Client, string) {
	err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

	return minioClient, bucketName
}

func uploadFile() {
	ctx := context.Background()
	minioClient, bucketName := GetClient()
	// Upload the test file
	// Change the value of filePath if the file is in another location
	objectName := "docker-compose.yml"
	filePath := "../docker-compose.yml"
	contentType := "application/octet-stream"

	// Upload the test file with FPutObject
	info, err := minioClient.FPutObject(
		ctx,
		bucketName,
		objectName,
		filePath,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(objectName, info)
}
