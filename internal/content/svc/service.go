package svc

import (
	"context"

	contentEntity "github.com/ghulammuzz/misterblast/internal/content/entity"
	contentRepo "github.com/ghulammuzz/misterblast/internal/content/repo"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

type ContentService interface {
	Add(content contentEntity.Content, lang string) error
	List(ctx context.Context, filter map[string]string) (*response.PaginateResponse, error)
	Delete(id int32) error
	Detail(ctx context.Context, id int32) (contentEntity.Content, error)
	Edit(id int32, content contentEntity.Content) error
}

type contentService struct {
	repo contentRepo.ContentRepository
}

func NewContentService(repo contentRepo.ContentRepository) ContentService {
	return &contentService{repo: repo}
}

func (s *contentService) Add(content contentEntity.Content, lang string) error {
	return s.repo.Add(content, lang)
}

func (s *contentService) List(ctx context.Context, filter map[string]string) (*response.PaginateResponse, error) {
	return s.repo.List(filter, ctx)
}

func (s *contentService) Delete(id int32) error {
	return s.repo.Delete(id)
}

func (s *contentService) Detail(ctx context.Context, id int32) (contentEntity.Content, error) {
	return s.repo.Detail(ctx, id)
}

func (s *contentService) Edit(id int32, content contentEntity.Content) error {
	return s.repo.Edit(id, content)
}
