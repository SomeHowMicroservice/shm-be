package common

import "gorm.io/gorm"

type Base64UploadRequest struct {
	Base64Data string `json:"base64_data"`
	FileName   string `json:"file_name"`
	Folder     string `json:"folder"`
}

type UploadFileResponse struct {
	FileID       string `json:"file_id"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Size         int64  `json:"size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
}

type Preload struct {
	Relation string
	Scope    func(*gorm.DB) *gorm.DB
}
