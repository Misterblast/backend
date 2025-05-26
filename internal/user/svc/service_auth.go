package svc

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	userEntity "github.com/ghulammuzz/misterblast/internal/user/entity"
	"github.com/ghulammuzz/misterblast/internal/user/repo"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
)

func (s *userService) Register(user userEntity.RegisterDTO) error {
	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	log.Info("[RegisterSvc] Start AddRepo")
	startAddRepo := time.Now()

	regUser := userEntity.Register{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	isVerified := true

	id, err := s.userRepo.Add(regUser, isVerified)
	if err != nil {
		log.Error("[UserSvc] Failed to register user", "error", err)
		return err
	}
	log.Info("[RegisterSvc] End AddRepo", "Total Duration", time.Since(startAddRepo))

	if user.Img != nil {
		go func(userImg *multipart.FileHeader, userID int64, svc *userService) {
			log.Info("[RegisterSvc] Start UploadImg")
			startUpload := time.Now()

			url, err := repo.ImageUploadProxyRESTY(userImg, fmt.Sprintf("/prod/user/profile-img/%d", userID))
			if err != nil {
				log.Error("[UserSvc] Failed to upload user image", "error", err)
				return
			}
			log.Info("[RegisterSvc] End UploadImg", "Total Duration", time.Since(startUpload))

			log.Info("[RegisterSvc] Start UpdateImg")
			startUpdate := time.Now()
			if err := svc.userRepo.UpdateImageURL(userID, url); err != nil {
				log.Error("[UserSvc] Failed to update user image URL", "error", err)
				return
			}
			log.Info("[RegisterSvc] End UpdateImg", "Total Duration", time.Since(startUpdate))
		}(user.Img, id, s)
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
