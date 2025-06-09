package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ghulammuzz/misterblast/internal/user/entity"
	"github.com/ghulammuzz/misterblast/internal/user/handler"
	"github.com/ghulammuzz/misterblast/pkg/middleware"
	"github.com/ghulammuzz/misterblast/pkg/response"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) SummaryUser(id int32, filter map[string]string) (*entity.UserSummary, error) {
	args := m.Called(id, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserSummary), args.Error(1)
}

func (m *MockUserService) Register(user entity.RegisterDTO) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) RegisterAdmin(admin entity.RegisterAdmin) error {
	args := m.Called(admin)
	return args.Error(0)
}

func (m *MockUserService) Login(user entity.UserLogin) (*entity.LoginResponse, string, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, "", args.Error(2)
	}
	return args.Get(0).(*entity.LoginResponse), args.String(1), args.Error(2)
}

func (m *MockUserService) AuthUser(userID int32) (entity.UserAuth, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return entity.UserAuth{}, args.Error(1)
	}
	return args.Get(0).(entity.UserAuth), args.Error(1)
}

func (m *MockUserService) DeleteUser(userID int32) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserService) DetailUser(userID int32) (entity.DetailUser, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return entity.DetailUser{}, args.Error(1)
	}
	return args.Get(0).(entity.DetailUser), args.Error(1)
}

func (m *MockUserService) EditUser(userID int32, user entity.EditUser) error {
	args := m.Called(userID, user)
	return args.Error(0)
}

func (m *MockUserService) ListUser(filters map[string]string, limit int, offset int) (*response.PaginateResponse, error) {
	args := m.Called(filters, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.PaginateResponse), args.Error(1)
}

// not implements
func (m *MockUserService) ChangePassword(token string, newPassword string) error {
	args := m.Called(token, newPassword)
	return args.Error(0)
}
func TestRegisterHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)
	h := handler.NewUserHandler(mockService, validator.New())
	app.Post("/register", h.RegisterHandler)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	_ = writer.WriteField("name", "John Doe")
	_ = writer.WriteField("email", "john@example.com")
	_ = writer.WriteField("password", "password")
	writer.Close()

	expected := entity.RegisterDTO{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password",
	}

	mockService.On("Register", expected).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/register", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestRegisterAdminHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)
	h := handler.NewUserHandler(mockService, validator.New())
	app.Post("/admin-check", h.RegisterAdminHandler)

	user := entity.RegisterAdmin{Name: "John Doe", Email: "john@example.com"}
	mockService.On("RegisterAdmin", user).Return(nil)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/admin-check", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestLoginHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)
	h := handler.NewUserHandler(mockService, validator.New())
	app.Post("/login", h.LoginHandler)

	user := entity.UserLogin{Email: "john@example.com", Password: "password"}
	userJWT := &entity.LoginResponse{ID: 1, Email: "john@example.com", IsAdmin: false, IsVerified: true}
	token := "valid_token"
	mockService.On("Login", user).Return(userJWT, token, nil)

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestDeleteUserHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)
	h := handler.NewUserHandler(mockService, validator.New())
	app.Delete("/users/:id", h.DeleteUserHandler)

	mockService.On("DeleteUser", int32(1)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestDetailUserHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)
	h := handler.NewUserHandler(mockService, validator.New())
	app.Get("/users/:id", h.DetailUserHandler)

	user := entity.DetailUser{ID: 1, Name: "John Doe", Email: "john@example.com"}
	mockService.On("DetailUser", int32(1)).Return(user, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestEditUserHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)
	h := handler.NewUserHandler(mockService, validator.New())
	app.Put("/users/:id", h.EditUserHandler)

	userEdit := entity.EditUser{Name: "John Updated", Email: "john@edit.com"}
	mockService.On("EditUser", int32(1), userEdit).Return(nil)

	body, _ := json.Marshal(userEdit)
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestListUserHandler(t *testing.T) {
	app := fiber.New()
	mockService := new(MockUserService)
	h := handler.NewUserHandler(mockService, validator.New())
	app.Get("/users", h.ListUsersHandler)

	mockUsers := []entity.ListUser{
		{ID: 1, Name: "John Doe", Email: "john@example.com", ImgUrl: ""},
	}

	mockResponse := &response.PaginateResponse{
		Total: int64(len(mockUsers)),
		Page:  1,
		Limit: 10,
		Data:  mockUsers,
	}

	mockService.On("ListUser", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode) // ✅ Hanya cek status code 200

	mockService.AssertExpectations(t)
}

func TestMeUserHandler(t *testing.T) {
	app := fiber.New()
	mockUserService := new(MockUserService)
	handler := handler.NewUserHandler(mockUserService, validator.New())
	app.Get("/me", middleware.JWTProtected(), handler.MeUserHandler)

	t.Run("Success - Valid Token", func(t *testing.T) {
		claims := jwt.MapClaims{
			"apps":     "misterblast-core",
			"email":    "john@example.com",
			"user_id":  1,
			"is_admin": false,
			"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		// fmt.Println(signedToken)

		mockUserService.On("AuthUser", int32(1)).Return(entity.UserAuth{
			ID:         1,
			Name:       "John Doe",
			Email:      "john@example.com",
			ImgUrl:     "http://example.com/image.jpg",
			IsAdmin:    false,
			IsVerified: true,
		}, nil)

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", "Bearer "+signedToken)
		resp, _ := app.Test(req)

		// body, _ := io.ReadAll(resp.Body)
		// fmt.Println(string(body))

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Fail - Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Fail - No Token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("Fail - User Not Found", func(t *testing.T) {
		claims := jwt.MapClaims{"user_id": float64(2)}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

		mockUserService.On("AuthUser", int32(2)).Return(entity.UserAuth{}, errors.New("user not found"))

		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", "Bearer "+signedToken)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
