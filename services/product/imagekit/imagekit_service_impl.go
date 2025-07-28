package imagekit

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	customErr "github.com/SomeHowMicroservice/shm-be/common/errors"
	"github.com/SomeHowMicroservice/shm-be/services/product/common"
	"github.com/SomeHowMicroservice/shm-be/services/product/config"
	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

type imageKitServiceImpl struct {
	client *imagekit.ImageKit
}

func NewImageKitService(cfg *config.Config) ImageKitService {
	client := imagekit.NewFromParams(imagekit.NewParams{
		PrivateKey:  cfg.ImageKit.PrivateKey,
		PublicKey:   cfg.ImageKit.PublicKey,
		UrlEndpoint: cfg.ImageKit.URLEndpoint,
	})

	return &imageKitServiceImpl{client}
}

func (s *imageKitServiceImpl) UploadFile(ctx context.Context, req *common.UploadFileRequest) (*common.UploadFileResponse, error) {
	if err := validateFile(req.FileName); err != nil {
		return nil, err
	}

	params := uploader.UploadParam{
		FileName: req.FileName,
		UseUniqueFileName: boolPtr(false),
	}

	if req.Folder != "" {
		params.Folder = req.Folder
	}

	result, err := s.client.Uploader.Upload(ctx, req.File, params)
	if err != nil {
		return nil, fmt.Errorf("upload file thất bại: %w", err)
	}

	return &common.UploadFileResponse{
		FileID:       result.Data.FileId,
		Name:         result.Data.Name,
		URL:          result.Data.Url,
		ThumbnailURL: result.Data.ThumbnailUrl,
		Size:         int64(result.Data.Size),
		Width:        result.Data.Width,
		Height:       result.Data.Height,
	}, nil
}

func (s *imageKitServiceImpl) UploadFromMultipart(ctx context.Context, fileHeader *multipart.FileHeader, folder string) (*common.UploadFileResponse, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("mở file thất bại: %w", err)
	}
	defer file.Close()

	req := &common.UploadFileRequest{
		File: file,
		FileName: fileHeader.Filename,
		Folder: folder,
	}

	return s.UploadFile(ctx, req)
}

func validateFile(fileName string) error {
	ext := strings.ToLower(filepath.Ext(fileName))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg"}

	for _, allowed := range allowedExts {
		if ext == allowed {
			return nil
		}
	}

	return customErr.ErrUnSupportedFileType
}

func boolPtr(b bool) *bool {
	return &b
}
