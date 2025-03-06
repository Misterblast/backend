package svc

import (
	"errors"
	"os"

	userEntity "github.com/ghulammuzz/misterblast/internal/user/entity"
)

func (s *userService) Register(user userEntity.Register) error {
	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return s.userRepo.Add(user, true)
}

func (s *userService) RegisterAdmin(user userEntity.RegisterAdmin) error {

	// check in csv or excel

	return s.userRepo.Add(userEntity.Register{
		Name:     user.Name,
		Email:    user.Email,
		Password: os.Getenv("PASSWORD_ALG"),
	}, false)
}
