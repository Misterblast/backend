package svc

import (
	setEntity "github.com/ghulammuzz/misterblast/internal/set/entity"
	setRepo "github.com/ghulammuzz/misterblast/internal/set/repo"
)

type SetService interface {
	AddSet(set setEntity.SetSet) error
	DeleteSet(id int32) error
	ListSets(filter map[string]string) ([]setEntity.ListSet, error)
}

type setService struct {
	repo setRepo.SetRepository
}

func NewSetService(repo setRepo.SetRepository) SetService {
	return &setService{repo: repo}
}

func (s *setService) AddSet(set setEntity.SetSet) error {
	return s.repo.Add(set)
}

func (s *setService) DeleteSet(id int32) error {
	return s.repo.Delete(id)
}

func (s *setService) ListSets(filter map[string]string) ([]setEntity.ListSet, error) {
	return s.repo.List(filter)
}
