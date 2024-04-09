package upload

import (
	"checkpoint/utils"
	"context"
	"log"

	"github.com/google/uuid"
)

type UploadResolver struct{}

func (UploadResolver) UploadFile(ctx context.Context, input UploadFileInput) (*UploadResult, error) {

	file, err := input.File.CreateReadStream()

	if err != nil {
		log.Println("upload-file-create-read-stream-error", err.Error())

		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	key, url, err := UploadFile(ctx, uuid.New().String()+":"+input.File.FileName, file)

	if err != nil {
		log.Println("upload-file-error", err.Error())

		return nil, utils.GraphqlError{
			Code:    err.Error(),
			Message: err.Error(),
		}
	}

	return &UploadResult{
		S3Key: *key,
		Url:   *url,
	}, nil
}

func (UploadResolver) UploadFiles(ctx context.Context, input UploadFilesInput) ([]UploadResult, error) {
	return []UploadResult{}, nil
}
