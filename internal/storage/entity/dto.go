package entity

import "mime/multipart"

type UploadFileRequestDto struct {
	Key  string                `form:"key"`
	File *multipart.FileHeader `form:"file"`
}

type UploadFileResponseDto struct {
	Url string `json:"url"`
}
