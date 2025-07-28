package common

import "io"

type UploadFileRequest struct {
	File     io.Reader
	FileName string
	Folder   string
}

type UploadFileResponse struct {
	FileID       string   `json:"file_id"`
	Name         string   `json:"name"`
	URL          string   `json:"url"`
	ThumbnailURL string   `json:"thumbnail_url"`
	Size         int64    `json:"size"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
}
