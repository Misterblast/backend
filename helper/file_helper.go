package helper

import (
	"fmt"
	"mime/multipart"
	"strings"
)

func GetFileType(file *multipart.FileHeader) (string, error) {
	filename := file.Filename
	ext := strings.ToLower(filename[strings.LastIndex(filename, "."):])

	imageExtensions := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
	}

	if imageExtensions[ext] {
		return "img", nil
	} else if ext == ".pdf" {
		return "pdf", nil
	}
	return "", fmt.Errorf("unsupported file type: %s", ext)
}

func ValidateFileSize(file *multipart.FileHeader, maxSize int64) bool {
	return file.Size <= maxSize
}
