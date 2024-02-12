package upload

import (
	"checkpoint/s3"
	"fmt"

	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/minio/minio-go/v7"
)

const maxSize = 8 * iris.MB

func UploadController(ctx iris.Context) {
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

	minioClient, err := s3.GetClient()

	if err != nil {
		ctx.JSON(iris.Map{"error": "Failed to upload file"})
	}

	// Upload the file to MinIO
	_, err = minioClient.Client.PutObject(ctx, minioClient.BucketName, objectName, file, info.Size, minio.PutObjectOptions{})

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to upload file"})
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

func GetSignedURLController(ctx iris.Context) {

	objectName := ctx.URLParam("name")

	if objectName == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid file object"})
	}

	// Generate a signed URL for the object
	presignedURL, err := signedUrl(ctx, objectName)

	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to get file"})
	}

	if presignedURL == nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"error": "file not found"})
	}

	ctx.JSON(iris.Map{
		"url": presignedURL,
	})

}
