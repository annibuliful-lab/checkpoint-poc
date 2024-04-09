package upload

import uploadmiddleware "checkpoint/gql/upload-middleware"

type UploadResult struct {
	S3Key string
	Url   string
}

type UploadFileInput struct {
	File uploadmiddleware.GraphQLUpload
}

type UploadFilesInput struct {
	Files []uploadmiddleware.GraphQLUpload
}
