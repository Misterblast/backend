package di

import (
	"database/sql"

	emailHandler "github.com/ghulammuzz/misterblast/internal/email/handler"
	emailRepo "github.com/ghulammuzz/misterblast/internal/email/repo"
	emailSvc "github.com/ghulammuzz/misterblast/internal/email/svc"
	userRepo "github.com/ghulammuzz/misterblast/internal/user/repo"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func InitializedEmailServiceFake(sb *sql.DB, val *validator.Validate) *emailHandler.EmailHandler {
	wire.Build(
		emailHandler.NewEmailHandler,
		emailSvc.NewEmailService,
		emailRepo.NewEmailRepository,
		emailRepo.NewOTPService,
		userRepo.NewUserRepository,
	)

	return &emailHandler.EmailHandler{}
}
