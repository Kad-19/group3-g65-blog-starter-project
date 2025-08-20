package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"g3-g65-bsp/domain"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) Register(ctx context.Context, email, username, password string) error {
	args := m.Called(ctx, email, username, password)
	return args.Error(0)
}

func (m *MockAuthUsecase) Login(ctx context.Context, email, password string) (string, string, int, *domain.User, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.String(1), args.Int(2), args.Get(3).(*domain.User), args.Error(4)
}

func (m *MockAuthUsecase) ActivateUser(ctx context.Context, token, email string) error {
	args := m.Called(ctx, token, email)
	return args.Error(0)
}

func (m *MockAuthUsecase) ResendActivationEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockAuthUsecase) ForgotPassword(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockAuthUsecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	args := m.Called(ctx, token, newPassword)
	return args.Error(0)
}

func (m *MockAuthUsecase) RefreshTokens(ctx context.Context, refreshToken string) (string, string, int, error) {
	args := m.Called(ctx, refreshToken)
	return args.String(0), args.String(1), args.Int(2), args.Error(3)
}

func (m *MockAuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockAuthUsecase) LogoutAll(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestAuthController_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockAuthUsecase := new(MockAuthUsecase)
		authController := NewAuthController(mockAuthUsecase, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := UserCreateRequest{Username: "testuser", Email: "test@example.com", Password: "password"}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockAuthUsecase.On("Register", mock.Anything, reqBody.Email, reqBody.Username, reqBody.Password).Return(nil)

		authController.Register(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockAuthUsecase.AssertExpectations(t)
	})

	t.Run("bad request", func(t *testing.T) {
		mockAuthUsecase := new(MockAuthUsecase)
		authController := NewAuthController(mockAuthUsecase, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte("")))
		c.Request.Header.Set("Content-Type", "application/json")

		authController.Register(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("conflict", func(t *testing.T) {
		mockAuthUsecase := new(MockAuthUsecase)
		authController := NewAuthController(mockAuthUsecase, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := UserCreateRequest{Username: "testuser", Email: "test@example.com", Password: "password"}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockAuthUsecase.On("Register", mock.Anything, reqBody.Email, reqBody.Username, reqBody.Password).Return(errors.New("user already exists"))

		authController.Register(c)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockAuthUsecase.AssertExpectations(t)
	})
}

func TestAuthController_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockAuthUsecase := new(MockAuthUsecase)
		authController := NewAuthController(mockAuthUsecase, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := UserLoginRequest{Email: "test@example.com", Password: "password"}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		user := &domain.User{ID: "1", Email: "test@example.com", Username: "testuser", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		mockAuthUsecase.On("Login", mock.Anything, reqBody.Email, reqBody.Password).Return("access_token", "refresh_token", 3600, user, nil)

		authController.Login(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockAuthUsecase.AssertExpectations(t)
	})

	t.Run("bad request", func(t *testing.T) {
		mockAuthUsecase := new(MockAuthUsecase)
		authController := NewAuthController(mockAuthUsecase, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("")))
		c.Request.Header.Set("Content-Type", "application/json")

		authController.Login(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockAuthUsecase := new(MockAuthUsecase)
		authController := NewAuthController(mockAuthUsecase, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		reqBody := UserLoginRequest{Email: "test@example.com", Password: "password"}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockAuthUsecase.On("Login", mock.Anything, reqBody.Email, reqBody.Password).Return("", "", 0, (*domain.User)(nil), errors.New("invalid credentials"))

		authController.Login(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockAuthUsecase.AssertExpectations(t)
	})
}
