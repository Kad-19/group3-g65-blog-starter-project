package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"g3-g65-bsp/domain"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Promote(ctx context.Context, adminID, userEmail string) error {
	args := m.Called(ctx, adminID, userEmail)
	return args.Error(0)
}

func (m *MockUserUsecase) Demote(ctx context.Context, adminID, userEmail string) error {
	args := m.Called(ctx, adminID, userEmail)
	return args.Error(0)
}

func (m *MockUserUsecase) ProfileUpdate(ctx context.Context, userID, bio, contactInfo string, profilePicture io.Reader) error {
	args := m.Called(ctx, userID, bio, contactInfo, profilePicture)
	return args.Error(0)
}

func (m *MockUserUsecase) GetAllUsers(ctx context.Context, page, limit int) ([]domain.User, int64, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]domain.User), args.Get(1).(int64), args.Error(2)
}

func TestUserController_ChangeUserRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("promote_success", func(t *testing.T) {
		mockUserUsecase := new(MockUserUsecase)
		userController := NewUserController(mockUserUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "admin123")

		reqBody := EmailReq{Email: "user@example.com"}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/users/promote", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockUserUsecase.On("Promote", mock.Anything, "admin123", "user@example.com").Return(nil)

		userController.HandlePromote(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("demote_success", func(t *testing.T) {
		mockUserUsecase := new(MockUserUsecase)
		userController := NewUserController(mockUserUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "admin123")

		reqBody := EmailReq{Email: "user@example.com"}
		jsonBody, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/users/demote", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		mockUserUsecase.On("Demote", mock.Anything, "admin123", "user@example.com").Return(nil)

		userController.HandleDemote(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockUserUsecase := new(MockUserUsecase)
		userController := NewUserController(mockUserUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/users/promote", nil)

		userController.HandlePromote(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestUserController_HandleUpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUserUsecase := new(MockUserUsecase)
		userController := NewUserController(mockUserUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("user_id", "user123")

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		writer.WriteField("bio", "new bio")
		writer.WriteField("contact_info", "new contact")
		part, _ := writer.CreateFormFile("profile_picture", "test.jpg")
		part.Write([]byte("test image"))
		writer.Close()

		c.Request, _ = http.NewRequest(http.MethodPut, "/users/profile", body)
		c.Request.Header.Set("Content-Type", writer.FormDataContentType())

		mockUserUsecase.On("ProfileUpdate", mock.Anything, "user123", "new bio", "new contact", mock.Anything).Return(nil)

		userController.HandleUpdateUser(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})
}

func TestUserController_HandleGetAllUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUserUsecase := new(MockUserUsecase)
		userController := NewUserController(mockUserUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/users?page=1&limit=10", nil)

		users := []domain.User{{ID: "1", Username: "test"}}
		mockUserUsecase.On("GetAllUsers", mock.Anything, 1, 10).Return(users, int64(1), nil)

		userController.HandleGetAllUsers(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotNil(t, resp["data"])
		mockUserUsecase.AssertExpectations(t)
	})

	t.Run("usecase_error", func(t *testing.T) {
		mockUserUsecase := new(MockUserUsecase)
		userController := NewUserController(mockUserUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/users", nil)

		mockUserUsecase.On("GetAllUsers", mock.Anything, 1, 10).Return([]domain.User{}, int64(0), errors.New("some error"))

		userController.HandleGetAllUsers(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUserUsecase.AssertExpectations(t)
	})
}
