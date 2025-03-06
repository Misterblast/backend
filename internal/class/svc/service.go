package svc

import (
	classEntity "github.com/ghulammuzz/misterblast/internal/class/entity"
	classRepo "github.com/ghulammuzz/misterblast/internal/class/repo"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
)

type ClassService interface {
	AddClass(class classEntity.SetClass) error
	DeleteClass(id int32) error
	ListClasses() ([]classEntity.Class, error)
}

type classService struct {
	repo classRepo.ClassRepository
}

func NewClassService(repo classRepo.ClassRepository) ClassService {
	return &classService{repo: repo}
}

func (s *classService) AddClass(class classEntity.SetClass) error {
	if class.Name == "" {
		log.Error("[Svc][AddClass] Error: name is required")
		return app.NewAppError(400, "name is required")
	}
	exists, err := s.repo.Exists(class.Name)
	if err != nil {
		log.Error("[Svc][AddLesson] Error: ", err)
		return err
	}

	if exists {
		log.Error("[Svc][AddLesson] Error: lesson already exists")
		return app.NewAppError(400, "lesson already exists")
	}

	err = s.repo.Add(class)
	if err != nil {
		log.Error("[Svc][AddClass] Error: ", err)
		return err
	}

	return nil
}

func (s *classService) DeleteClass(id int32) error {
	if id <= 0 {
		log.Error("[Svc][DeleteClass] Error: invalid id")
		return app.NewAppError(400, "invalid id")
	}

	err := s.repo.Delete(id)
	if err != nil {
		log.Error("[Svc][DeleteClass] Error: ", err)
		return err
	}

	return nil
}

func (s *classService) ListClasses() ([]classEntity.Class, error) {
	classes, err := s.repo.List()
	if err != nil {
		log.Error("[Svc][ListClasses] Error: ", err)
		return nil, err
	}

	return classes, nil
}
