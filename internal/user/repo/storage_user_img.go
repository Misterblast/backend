package repo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/ghulammuzz/misterblast/internal/user/entity"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
)

func ImageUploadProxy(fileHeader *multipart.FileHeader, key string) (string, error) {
	// Validasi input
	if fileHeader == nil {
		return "", errors.New("file header is nil")
	}
	if key == "" {
		return "", errors.New("key is empty")
	}

	token := os.Getenv("FILE_SERVICE_TOKEN")
	if token == "" {
		return "", errors.New("FILE_SERVICE_TOKEN environment variable not set")
	}

	// Prepare the multipart form
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add the key field
	if err := writer.WriteField("key", key); err != nil {
		log.Error("ImageUploadProxy", "Failed to write field 'key'", err)
		return "", err
	}

	// Add the file
	fw, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		log.Error("ImageUploadProxy", "Failed to create form file", err)
		return "", err
	}

	src, err := fileHeader.Open()
	if err != nil {
		log.Error("ImageUploadProxy", "Failed to open file header", err)
		return "", err
	}
	defer src.Close()

	if _, err := io.Copy(fw, src); err != nil {
		log.Error("ImageUploadProxy", "Failed to copy file content", err)
		return "", err
	}

	// Close the writer before making the request
	writer.Close()

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create request
	req, err := http.NewRequest("POST", "https://stg.file.go-assessment.link/file/", &buf)
	if err != nil {
		log.Error("ImageUploadProxy", "Failed to create request", err)
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("MISTERBLAST_API_KEY", token)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		log.Error("ImageUploadProxy", "Failed to upload image", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error("ImageUploadProxy", "Upload failed with status code", resp.StatusCode)
		return "", fmt.Errorf("upload failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("ImageUploadProxy", "Failed to read response body", err)
		return "", err
	}

	if len(body) == 0 {
		return "", errors.New("empty response body")
	}

	// Parse response
	var response entity.ImgResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Error("ImageUploadProxy", "Failed to unmarshal response", err)
		return "", err
	}

	if response.Data.URL == "" {
		log.Error("ImageUploadProxy", "No URL returned from server")
		return "", errors.New("no URL returned from server")
	}

	log.Info("ImageUploadProxy", "Image uploaded successfully", "URL", response.Data.URL)
	return response.Data.URL, nil
}
