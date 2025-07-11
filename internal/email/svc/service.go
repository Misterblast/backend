package svc

import (
	"time"

	emailRepo "github.com/ghulammuzz/misterblast/internal/email/repo"
	userRepo "github.com/ghulammuzz/misterblast/internal/user/repo"
	"github.com/ghulammuzz/misterblast/pkg/app"
	log "github.com/ghulammuzz/misterblast/pkg/middleware"
)

type EmailService interface {
	SendOTP(email string) error
	SendDeeplink(email string) (string, error)
	ValidateOTP(adminID int32, otp string) error
}

func NewEmailService(emailRepo emailRepo.EmailRepository, userRepo userRepo.UserRepository, otp emailRepo.OTP) EmailService {
	return &emailService{
		emailRepo: emailRepo,
		userRepo:  userRepo,
		otp:       otp,
	}
}

type emailService struct {
	emailRepo emailRepo.EmailRepository
	userRepo  userRepo.UserRepository
	otp       emailRepo.OTP
}

func (s *emailService) SendOTP(email string) error {

	adminID, err := s.userRepo.GetIDByEmail(email)
	if err != nil {
		return err
	}
	otpString, err := s.otp.GenerateOTP()
	if err != nil {
		return err
	}

	expAt := time.Now().Add(120 * time.Second).Unix()

	if err := s.emailRepo.SetOTP(adminID, otpString, expAt); err != nil {
		return err
	}

	if err := s.otp.SendEmailSMTP(email, otpString); err != nil {
		return err
	}

	return nil
}

func (s *emailService) ValidateOTP(adminID int32, otp string) error {

	exists, err := s.userRepo.Exists(adminID)
	if !exists {
		if err != nil {
			log.Error("[Svc][userRepo.Exists] Error Exec: ", err)
			return app.NewAppError(500, err.Error())
		}
	}

	dbOtp, expiresAt, err := s.emailRepo.GetOTP(adminID)
	if err != nil {
		log.Error("[Svc][s.emailRepo.GetOTP] Error Exec: ", err)
		return app.NewAppError(500, err.Error())
	}

	if dbOtp != otp {
		log.Error("[Svc][otp] Error Exec: ", err)
		return app.NewAppError(400, "OTP tidak sesuai")
	}

	if time.Now().Unix() > expiresAt {
		log.Error("[Svc][expiresAt] Error Exec: ", err)
		return app.NewAppError(400, "OTP sudah kedaluwarsa")
	}

	err = s.userRepo.AdminActivation(adminID)
	if err != nil {
		log.Error("[Svc][s.userRepo.AdminActivation] Error Exec: ", err)
		return app.NewAppError(500, err.Error())
	}

	return nil
}

func (s *emailService) SendDeeplink(email string) (string, error) {

	userID, err := s.userRepo.GetIDByEmail(email)
	if err != nil {
		log.Error("[Svc][s.userRepo.GetIDByEmail] Error Exec: ", err)
		return "", err
	}
	tokenString, err := s.userRepo.GenerateToken()
	if err != nil {
		log.Error("[Svc][s.userRepo.GenerateToken] Error Exec: ", err)
		return "", err
	}

	expAt := time.Now().Add(120 * time.Second).Unix()

	if err := s.userRepo.SetDeeplink(userID, tokenString, expAt); err != nil {
		log.Error("[Svc][s.userRepo.SetDeeplink] Error Exec: ", err)
		return "", err
	}

	if err := s.otp.SendDeeplinkEmailSMTP(email, tokenString); err != nil {
		log.Error("[Svc][s.otp.SendDeeplinkEmailSMTP] Error Exec: ", err)
		return "", err
	}

	return tokenString, nil
}
