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
	"github.com/go-resty/resty/v2"
)

func ImageUploadProxyHTTP(fileHeader *multipart.FileHeader, key string) (string, error) {
	start := time.Now()
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

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	log.Debug("ImageUploadProxy", "start writeField")
	if err := writer.WriteField("key", key); err != nil {
		log.Error("ImageUploadProxy", "Failed to write field 'key'", err)
		return "", err
	}
	log.Debug("ImageUploadProxy", "writeField duration:", time.Since(start))

	start = time.Now()
	fw, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		log.Error("ImageUploadProxy", "Failed to create form file", err)
		return "", err
	}
	log.Debug("ImageUploadProxy", "io.Copy duration:", time.Since(start))

	src, err := fileHeader.Open()
	if err != nil {
		log.Error("ImageUploadProxy", "Failed to open file header", err)
		return "", err
	}

	defer src.Close()

	start = time.Now()
	if _, err := io.Copy(fw, src); err != nil {
		log.Error("ImageUploadProxy", "Failed to copy file content", err)
		return "", err
	}
	log.Debug("ImageUploadProxy", "io.Copy duration:", time.Since(start))

	writer.Close()

	start = time.Now()
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", "https://stg.file.go-assessment.link/file/", &buf)
	if err != nil {
		log.Error("ImageUploadProxy", "Failed to create request", err)
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("MISTERBLAST_API_KEY", token)

	resp, err := client.Do(req)
	if err != nil {
		log.Error("ImageUploadProxy", "Failed to upload image", err)
		return "", err
	}
	log.Debug("ImageUploadProxy", "HTTP request duration:", time.Since(start))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error("ImageUploadProxy", "Upload failed with status code", resp.StatusCode)
		return "", fmt.Errorf("upload failed with status code %d: %s", resp.StatusCode, string(body))
	}

	start = time.Now()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("ImageUploadProxy", "Failed to read response body", err)
		return "", err
	}

	if len(body) == 0 {
		return "", errors.New("empty response body")
	}
	log.Debug("ImageUploadProxy", "Read response body duration:", time.Since(start))

	start = time.Now()
	var response entity.ImgResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Error("ImageUploadProxy", "Failed to unmarshal response", err)
		return "", err
	}
	log.Debug("ImageUploadProxy", "Unmarshal response duration:", time.Since(start))

	if response.Data.URL == "" {
		log.Error("ImageUploadProxy", "No URL returned from server")
		return "", errors.New("no URL returned from server")
	}

	log.Debug("ImageUploadProxy", "Image uploaded successfully", "URL", response.Data.URL)
	return response.Data.URL, nil
}

func ImageUploadProxyRESTY(fileHeader *multipart.FileHeader, key string) (string, error) {
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

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if err := writer.WriteField("key", key); err != nil {
		return "", err
	}

	fw, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		return "", err
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	if _, err := io.Copy(fw, src); err != nil {
		return "", err
	}
	writer.Close()

	client := resty.New().
		SetTimeout(30 * time.Second)

	resp, err := client.R().
		SetHeader("Content-Type", writer.FormDataContentType()).
		SetHeader("MISTERBLAST_API_KEY", token).
		SetBody(buf.Bytes()).
		Post(os.Getenv("STORAGE_URL"))
	if err != nil {
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("upload failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	var response entity.ImgResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return "", err
	}

	if response.Data.URL == "" {
		return "", errors.New("no URL returned from server")
	}

	return response.Data.URL, nil
}
