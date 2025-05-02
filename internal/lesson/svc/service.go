package svc

import (
	"github.com/ghulammuzz/misterblast/internal/lesson/entity"
	"github.com/ghulammuzz/misterblast/internal/lesson/repo"
	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
)

type LessonService interface {
	AddLesson(lesson entity.Lesson) error
	DeleteLesson(id int32) error
	ListLessons() ([]entity.Lesson, error)
}

type lessonService struct {
	repo repo.LessonRepository
}

func NewLessonService(repo repo.LessonRepository) LessonService {
	return &lessonService{repo: repo}
}

func (s *lessonService) AddLesson(lesson entity.Lesson) error {
	if lesson.Name == "" {
		log.Error("[Svc][AddLesson] Error: name is required")
		return app.NewAppError(400, "name is required")
	}

	exists, err := s.repo.Exists(lesson.Name)
	if err != nil {
		log.Error("[Svc][AddLesson] Error: ", err)
		return err
	}

	if exists {
		log.Error("[Svc][AddLesson] Error: lesson already exists")
		return app.NewAppError(400, "lesson already exists")
	}

	err = s.repo.Add(lesson)
	if err != nil {
		log.Error("[Svc][AddLesson] Error: ", err)
		return err
	}

	return nil
}

func (s *lessonService) DeleteLesson(id int32) error {
	if id <= 0 {
		log.Error("[Svc][DeleteLesson] Error: invalid id")
		return app.NewAppError(400, "invalid id")
	}

	err := s.repo.Delete(id)
	if err != nil {
		log.Error("[Svc][DeleteLesson] Error: ", err)
		return err
	}

	return nil
}

func (s *lessonService) ListLessons() ([]entity.Lesson, error) {
	lessons, err := s.repo.List()
	if err != nil {
		log.Error("[Svc][ListLessons] Error: ", err)
		return nil, err
	}

	return lessons, nil
}
