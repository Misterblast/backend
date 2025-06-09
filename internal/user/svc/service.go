package svc

import (
	tQuizRepo "github.com/ghulammuzz/misterblast/internal/quiz/repo"
	tTaskRepo "github.com/ghulammuzz/misterblast/internal/task/repo"
	userEntity "github.com/ghulammuzz/misterblast/internal/user/entity"
	userRepo "github.com/ghulammuzz/misterblast/internal/user/repo"
	"github.com/ghulammuzz/misterblast/pkg/jwt"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

type UserService interface {
	Register(user userEntity.RegisterDTO) error
	RegisterAdmin(user userEntity.RegisterAdmin) error
	Login(user userEntity.UserLogin) (*userEntity.LoginResponse, string, error)
	ListUser(filter map[string]string, page, limit int) (*response.PaginateResponse, error)
	DetailUser(id int32) (userEntity.DetailUser, error)
	AuthUser(id int32) (userEntity.UserAuth, error)
	EditUser(id int32, user userEntity.EditUser) error
	DeleteUser(id int32) error
	ChangePassword(token string, newPassword string) error
	SummaryUser(id int32, filter map[string]string) (*userEntity.UserSummary, error)
}
type userService struct {
	userRepo  userRepo.UserRepository
	tQuizRepo tQuizRepo.QuizRepository
	tTaskRepo tTaskRepo.TaskRepository
}

func NewUserService(userRepo userRepo.UserRepository, tQuizRepo tQuizRepo.QuizRepository, tTaskRepo tTaskRepo.TaskRepository) UserService {
	return &userService{userRepo: userRepo, tQuizRepo: tQuizRepo, tTaskRepo: tTaskRepo}
}

func (s *userService) SummaryUser(id int32, filter map[string]string) (*userEntity.UserSummary, error) {
	quizCount, avgQuiz, err := s.tQuizRepo.GetAvgTotal(int(id), filter)
	if err != nil {
		return nil, err
	}

	taskCount, avgTask, err := s.tTaskRepo.GetAvgTotal(id)
	if err != nil {
		return nil, err
	}

	return &userEntity.UserSummary{
		TotalQuizAttempts: quizCount,
		AvgQuizScore:      avgQuiz,
		TotalTaskAttempts: taskCount,
		AvgTaskScore:      avgTask,
	}, nil
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

func (s *userService) ListUser(filter map[string]string, page, limit int) (*response.PaginateResponse, error) {
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

func (s *userService) ChangePassword(token string, newPassword string) error {

	deeplink, err := s.userRepo.GetDeeplink(token)
	if err != nil {
		return err
	}

	return s.userRepo.EditPassword(deeplink.UserID, newPassword)
}
