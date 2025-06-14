package svc

import (
	"context"

	"github.com/ghulammuzz/misterblast/internal/content/entity"
	"github.com/ghulammuzz/misterblast/internal/content/repo"
	"github.com/ghulammuzz/misterblast/pkg/app"
)

type AuthorService interface {
	CreateAuthor(ctx context.Context, req entity.CreateAuthorRequest) error
	UpdateAuthor(ctx context.Context, id int, req entity.UpdateAuthorRequest) error
	DeleteAuthor(ctx context.Context, id int32) error
	GetAuthor(ctx context.Context, id int32) (*entity.Author, error)
	ListAuthors(ctx context.Context) ([]entity.Author, error)
}

type authorService struct {
	repo repo.AuthorRepository
}

func NewAuthorService(r repo.AuthorRepository) AuthorService {
	return &authorService{r}
}

func (s *authorService) CreateAuthor(ctx context.Context, req entity.CreateAuthorRequest) error {
	if ok, _ := s.repo.Exists(req.Name); ok {
		return app.NewAppError(400, "author already exists")
	}
	return s.repo.Add(entity.Author{
		Name:        req.Name,
		ImgURL:      req.ImgURL,
		Description: req.Description,
	})
}

func (s *authorService) UpdateAuthor(ctx context.Context, id int, req entity.UpdateAuthorRequest) error {
	return s.repo.Update(entity.Author{
		ID:          id,
		Name:        req.Name,
		ImgURL:      req.ImgURL,
		Description: req.Description,
	})
}

func (s *authorService) DeleteAuthor(ctx context.Context, id int32) error {
	return s.repo.Delete(id)
}

func (s *authorService) GetAuthor(ctx context.Context, id int32) (*entity.Author, error) {
	return s.repo.Get(ctx, id)
}

func (s *authorService) ListAuthors(ctx context.Context) ([]entity.Author, error) {
	return s.repo.List(ctx)
}
