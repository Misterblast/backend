package svc

import (
	"errors"
	"fmt"
	"os"

	userEntity "github.com/ghulammuzz/misterblast/internal/user/entity"
	"github.com/ghulammuzz/misterblast/internal/user/repo"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
)

func (s *userService) Register(user userEntity.RegisterDTO) error {
	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	var regUser userEntity.Register
	regUser.Name = user.Name
	regUser.Email = user.Email
	regUser.Password = user.Password

	isVerified := true

	id, err := s.userRepo.Add(regUser, isVerified)
	if err != nil {
		log.Error("[UserSvc] Failed to register user", "error", err)
		return err
	}

	if user.Img != nil {

		url, err := repo.ImageUploadProxy(user.Img, fmt.Sprintf("/prod/user/profile-img/%d", id))
		if err != nil {
			log.Error("[UserSvc] Failed to upload user image", "error", err)
			return err
		}

		if err := s.userRepo.UpdateImageURL(id, url); err != nil {
			log.Error("[UserSvc] Failed to update user image URL", "error", err)
			return err
		}
	}

	return nil
}

func (s *userService) RegisterAdmin(user userEntity.RegisterAdmin) error {
	// TODO: implement check in csv/excel jika diperlukan

	// Pakai password dari env var (hardcoded?)
	regUser := userEntity.Register{
		Name:     user.Name,
		Email:    user.Email,
		Password: os.Getenv("PASSWORD_ALG"),
	}

	// Insert user, IsVerified = false
	_, err := s.userRepo.Add(regUser, false)
	if err != nil {
		return err
	}

	return nil
}
