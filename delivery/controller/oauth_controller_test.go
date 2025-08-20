package controller

import (
	"context"
	"encoding/json"
	"errors"
	"g3-g65-bsp/config"
	"g3-g65-bsp/domain"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type MockOAuthUsecase struct {
	mock.Mock
}

func (m *MockOAuthUsecase) OAuthLogin(ctx context.Context, cfg oauth2.Config, code string) (string, string, int, *domain.User, error) {
	args := m.Called(ctx, cfg, code)
	var user *domain.User
	if args.Get(3) != nil {
		user = args.Get(3).(*domain.User)
	}
	return args.String(0), args.String(1), args.Int(2), user, args.Error(4)
}

func TestOAuthController_HandleGoogleLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockOAuthUsecase := new(MockOAuthUsecase)
	oauthController := NewOAuthController(mockOAuthUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/login", nil)

	t.Setenv("OAUTH_STATE_STRING", "test-state")
	t.Setenv("GOOGLE_OAUTH_CLIENT_ID", "test-client-id")
	t.Setenv("GOOGLE_OAUTH_CLIENT_SECRET", "test-client-secret")
	t.Setenv("ACCESS_TOKEN_EXPIRY", "1m")
	t.Setenv("REFRESH_TOKEN_EXPIRY", "1m")
	config.LoadConfig()

	oauthController.HandleGoogleLogin(c)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	assert.Contains(t, w.Header().Get("Location"), "accounts.google.com/o/oauth2/auth")
}

func TestOAuthController_HandleGoogleCallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Setenv("OAUTH_STATE_STRING", "test-state")
	t.Setenv("GOOGLE_OAUTH_CLIENT_ID", "test-client-id")
	t.Setenv("GOOGLE_OAUTH_CLIENT_SECRET", "test-client-secret")
	t.Setenv("ACCESS_TOKEN_EXPIRY", "1m")
	t.Setenv("REFRESH_TOKEN_EXPIRY", "1m")
	config.LoadConfig()

	t.Run("success", func(t *testing.T) {
		mockOAuthUsecase := new(MockOAuthUsecase)
		oauthController := NewOAuthController(mockOAuthUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback?state=test-state&code=test-code", nil)

		user := &domain.User{ID: "user123", Email: "test@example.com"}
		mockOAuthUsecase.On("OAuthLogin", mock.Anything, mock.AnythingOfType("oauth2.Config"), "test-code").Return("access-token", "refresh-token", 3600, user, nil)

		oauthController.HandleGoogleCallback(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "access-token", resp["access_token"])
		assert.Equal(t, "refresh-token", resp["refresh_token"])
		mockOAuthUsecase.AssertExpectations(t)
	})

	t.Run("invalid_state", func(t *testing.T) {
		mockOAuthUsecase := new(MockOAuthUsecase)
		oauthController := NewOAuthController(mockOAuthUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback?state=invalid-state&code=test-code", nil)

		oauthController.HandleGoogleCallback(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("no_code", func(t *testing.T) {
		mockOAuthUsecase := new(MockOAuthUsecase)
		oauthController := NewOAuthController(mockOAuthUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback?state=test-state", nil)

		oauthController.HandleGoogleCallback(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase_error", func(t *testing.T) {
		mockOAuthUsecase := new(MockOAuthUsecase)
		oauthController := NewOAuthController(mockOAuthUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/auth/google/callback?state=test-state&code=test-code", nil)

		mockOAuthUsecase.On("OAuthLogin", mock.Anything, mock.AnythingOfType("oauth2.Config"), "test-code").Return("", "", 0, (*domain.User)(nil), errors.New("usecase error"))

		oauthController.HandleGoogleCallback(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockOAuthUsecase.AssertExpectations(t)
	})
}

// Mock implementation for domain.User for testing purposes
func (m *MockOAuthUsecase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockOAuthUsecase) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockOAuthUsecase) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockOAuthUsecase) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockOAuthUsecase) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOAuthUsecase) SetUserRole(ctx context.Context, id, role string) error {
	args := m.Called(ctx, id, role)
	return args.Error(0)
}

func (m *MockOAuthUsecase) SetUserLastLogin(ctx context.Context, id string, lastLogin time.Time) error {
	args := m.Called(ctx, id, lastLogin)
	return args.Error(0)
}
