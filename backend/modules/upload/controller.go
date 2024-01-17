package upload

import (
	"checkpoint/s3"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/minio/minio-go/v7"
)

const maxSize = 8 * iris.MB

func Upload(ctx iris.Context) {
	ctx.SetMaxRequestBodySize(maxSize)

	file, info, err := ctx.FormFile("file")

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to get the file from the form"})
		return
	}

	defer file.Close()

	// Create a unique object name based on the file name
	objectName := fmt.Sprintf("%s:%s", uuid.New(), info.Filename)

	minioClient, bucketName, err := s3.GetClient()

	if err != nil {
		ctx.JSON(iris.Map{"error": "Failed to upload file to MinIO"})
	}

	// Upload the file to MinIO
	_, err = minioClient.PutObject(ctx, bucketName, objectName, file, info.Size, minio.PutObjectOptions{})

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to upload file to MinIO"})
		return
	}
	url, err := signedUrl(ctx, objectName)

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)

		ctx.JSON(iris.Map{
			"error": "Failed to get file url",
		})
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{
		"message": "Upload compledted",
		"file":    objectName,
		"url":     url,
	})
}

func GetSignedURL(ctx iris.Context) {

	objectName := ctx.PostValue("objectName")

	if objectName == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid file"})
	}

	minioClient, bucketName, err := s3.GetClient()

	if err != nil {
		ctx.JSON(iris.Map{"error": "Failed to contect Storage"})
		return
	}

	// Generate a signed URL for the object
	presignedURL, err := minioClient.PresignedPutObject(ctx, bucketName, objectName, 15*time.Minute)

	if err != nil {
		ctx.JSON(iris.Map{"error": "Failed to get file"})

	}

	ctx.JSON(iris.Map{
		"url": presignedURL.String(),
	})

}
