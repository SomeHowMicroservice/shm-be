package imagekit

import (
	"context"
	"mime/multipart"

	"github.com/SomeHowMicroservice/shm-be/services/product/common"
)

type ImageKitService interface {
	UploadFile(ctx context.Context, req *common.UploadFileRequest) (*common.UploadFileResponse, error)

	UploadFromMultipart(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (*common.UploadFileResponse, error)
}