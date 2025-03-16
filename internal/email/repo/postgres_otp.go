package repo

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"os"

	"github.com/ghulammuzz/misterblast/pkg/app"
	"github.com/ghulammuzz/misterblast/pkg/log"
)

type OTP interface {
	GenerateOTP() (string, error)
	SendEmailSMTP(to string, otp string) error
	SendDeeplinkEmailSMTP(to string, deeplink string) error
}

type otpService struct{}

func NewOTPService() OTP {
	return &otpService{}
}

func (o *otpService) GenerateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", app.NewAppError(500, err.Error())
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func (o *otpService) SendEmailSMTP(to string, otp string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUser := os.Getenv("EMAIL_HOST_USER")
	smtpPassword := os.Getenv("EMAIL_HOST_PASSWORD")

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: OTP Verification\r\n" +
		"\r\n" +
		"Your OTP code is: " + otp + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, msg)
	if err != nil {
		log.Error("[Repo][smtp.SendMail.OTP] Error Exec: ", err)
		return app.NewAppError(500, err.Error())
	}
	return nil
}

func (o *otpService) SendDeeplinkEmailSMTP(to string, deeplink string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	smtpUser := os.Getenv("EMAIL_HOST_USER")
	smtpPassword := os.Getenv("EMAIL_HOST_PASSWORD")

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: OTP Verification\r\n" +
		"\r\n" +
		"Klik Tautan untuk mereset Password: " + deeplink + "\r\n")

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{to}, msg)
	if err != nil {
		log.Error("[Repo][smtp.SendMail.Deeplink] Error Exec: ", err)
		return app.NewAppError(500, err.Error())
	}
	return nil
}
