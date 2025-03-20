package helper

import "mime/multipart"

func GetFileExtension(file *multipart.FileHeader) string {
	return file.Filename[len(file.Filename)-4:]
}

func ValidateFileSize(file *multipart.FileHeader, maxSize int64) bool {
	return file.Size <= maxSize
}
