package svc

import (
	"errors"
	"io"
	"os"

	"github.com/ghulammuzz/misterblast/internal/models"
	"github.com/ghulammuzz/misterblast/internal/storage/entity"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type StorageService interface {
	UploadFile(request entity.UploadFileRequestDto) (models.Response[entity.UploadFileResponseDto], error)
	DeleteFile(key string) error
}

type StorageServiceImpl struct {
	client  *fiber.Client
	baseUrl string
	apiKey  string
}

func NewStorageService(client *fiber.Client) StorageService {
	baseUrl := os.Getenv("STORAGE_BASE_URL")
	apiKey := os.Getenv("STORAGE_API_KEY")
	return &StorageServiceImpl{client: client, baseUrl: baseUrl, apiKey: apiKey}
}

func (s StorageServiceImpl) UploadFile(request entity.UploadFileRequestDto) (models.Response[entity.UploadFileResponseDto], error) {
	file, err := request.File.Open()
	if err != nil {
		log.Errorf("[Svc][Storage]Failed to open file: %v", err)
		return models.Response[entity.UploadFileResponseDto]{}, app.NewAppError(fiber.StatusInternalServerError, "Failed to open file")
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Errorf("[Svc][Storage]Failed to read file: %v", err)
		return models.Response[entity.UploadFileResponseDto]{}, app.NewAppError(fiber.StatusInternalServerError, "Failed to read file")
	}

	formFile := fiber.FormFile{
		Fieldname: "file",
		Name:      request.File.Filename,
		Content:   fileBytes,
	}
	args := fiber.AcquireArgs()
	args.Set("key", request.Key)
	defer fiber.ReleaseArgs(args)

	req, err := s.prepareRequest(s.client.Post(s.baseUrl + "/img").FileData(&formFile).MultipartForm(args))
	if err != nil {
		var appErr *app.AppError
		if !errors.As(err, &appErr) {
			log.Errorf("[Svc][Storage]prepare request failed: %v", err)
			return models.Response[entity.UploadFileResponseDto]{}, app.NewAppError(fiber.StatusInternalServerError, "failed to prepare storage connection")
		}
		return models.Response[entity.UploadFileResponseDto]{}, appErr
	}
	var responseDto models.Response[entity.UploadFileResponseDto]
	code, _, errs := req.Struct(&responseDto)
	if len(errs) > 0 {
		log.Errorf("[Svc][Storage]Http request failed: %v", errs)
		return models.Response[entity.UploadFileResponseDto]{}, app.NewAppError(fiber.StatusInternalServerError, "Failed to upload file")
	}
	if code != fiber.StatusOK {
		log.Errorf("[Svc][Storage]Http request not ok: %v", responseDto)
		return models.Response[entity.UploadFileResponseDto]{}, app.NewAppError(fiber.StatusInternalServerError, "Failed to upload file")
	}
	return responseDto, nil
}

func (s StorageServiceImpl) DeleteFile(key string) error {
	//TODO implement me
	panic("implement me")
}

func (s StorageServiceImpl) prepareRequest(requestAgent *fiber.Agent) (*fiber.Agent, error) {
	if s.apiKey == "" {
		return nil, app.NewAppError(fiber.StatusInternalServerError, "Missing storage credentials")
	}

	return requestAgent.Set("MISTERBLAST_API_KEY", s.apiKey), nil
}
