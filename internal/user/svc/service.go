package svc

import (
	userEntity "github.com/ghulammuzz/misterblast/internal/user/entity"
	userRepo "github.com/ghulammuzz/misterblast/internal/user/repo"
	"github.com/ghulammuzz/misterblast/pkg/jwt"
)

type UserService interface {
	Register(user userEntity.Register) error
	RegisterAdmin(user userEntity.RegisterAdmin) error
	Login(user userEntity.UserLogin) (*userEntity.LoginResponse, string, error)
	ListUser(filter map[string]string, page, limit int) ([]userEntity.ListUser, error)
	DetailUser(id int32) (userEntity.DetailUser, error)
	AuthUser(id int32) (userEntity.UserAuth, error)
	EditUser(id int32, user userEntity.EditUser) error
	DeleteUser(id int32) error
}
type userService struct {
	userRepo userRepo.UserRepository
}

func NewUserService(userRepo userRepo.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Login(user userEntity.UserLogin) (*userEntity.LoginResponse, string, error) {

	var userResponse userEntity.LoginResponse

	userResult, err := s.userRepo.Check(user)
	if err != nil {
		return nil, "", err
	}

	userResponse.ID = userResult.ID
	userResponse.Email = userResult.Email
	userResponse.IsAdmin = userResult.IsAdmin
	userResponse.IsVerified = userResult.IsVerified

	token, err := jwt.GenerateJWT(*userResult)
	if err != nil {
		return nil, "", err
	}

	return &userResponse, token, nil
}

func (s *userService) ListUser(filter map[string]string, page, limit int) ([]userEntity.ListUser, error) {
	return s.userRepo.List(filter, page, limit)
}

func (s *userService) DetailUser(id int32) (userEntity.DetailUser, error) {
	return s.userRepo.Detail(id)
}

func (s *userService) EditUser(id int32, user userEntity.EditUser) error {
	return s.userRepo.Edit(id, user)
}

func (s *userService) AuthUser(id int32) (userEntity.UserAuth, error) {
	return s.userRepo.Auth(id)
}

func (s *userService) DeleteUser(id int32) error {
	return s.userRepo.Delete(id)
}
