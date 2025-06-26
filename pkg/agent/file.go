package agent

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

func FileUploadProxyRESTY(fileHeader *multipart.FileHeader, key string) (string, error) {
	
	if fileHeader == nil {
		log.Warn("[FileUploadProxyRESTY] file header is nil")
		return "", errors.New("file header is nil")
	}

	if key == "" {
		log.Warn("[FileUploadProxyRESTY] key is empty")
		return "", errors.New("key is empty")
	}

	token := os.Getenv("FILE_SERVICE_TOKEN")
	if token == "" {
		log.Warn("[FileUploadProxyRESTY] FILE_SERVICE_TOKEN environment variable not set")
		return "", errors.New("FILE_SERVICE_TOKEN environment variable not set")
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if err := writer.WriteField("key", key); err != nil {
		log.Error("[FileUploadProxyRESTY] failed to write field 'key'", "error", err)
		return "", err
	}

	fw, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		log.Error("[FileUploadProxyRESTY] failed to create form file", "error", err)
		return "", err
	}

	src, err := fileHeader.Open()
	if err != nil {
		log.Error("[FileUploadProxyRESTY] failed to open file header", "error", err)
		return "", err
	}
	defer src.Close()

	if _, err := io.Copy(fw, src); err != nil {
		log.Error("[FileUploadProxyRESTY] failed to copy file content", "error", err)
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
		log.Error("[FileUploadProxyRESTY] request failed", "error", err)
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Error("[FileUploadProxyRESTY] upload failed", "status_code", resp.StatusCode(), "body", string(resp.Body()))
		return "", fmt.Errorf("upload failed with status code %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	var response entity.ImgResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		log.Error("[FileUploadProxyRESTY] failed to unmarshal response", "error", err)
		return "", err
	}

	if response.Data.URL == "" {
		log.Error("[FileUploadProxyRESTY] no URL returned from server")
		return "", errors.New("no URL returned from server")
	}

	return response.Data.URL, nil
}
