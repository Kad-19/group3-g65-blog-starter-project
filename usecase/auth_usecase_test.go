package usecase

import (
	"context"
	"errors"
	"g3-g65-bsp/domain"
	"g3-g65-bsp/infrastructure/auth"
	"g3-g65-bsp/infrastructure/email"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockUserRepository mocks domain.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUserProfile(ctx context.Context, bio string, contactInfo string, imagePath string, Email string) error {
	args := m.Called(ctx, bio, contactInfo, imagePath, Email)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserRole(ctx context.Context, role string, Email string) error {
	args := m.Called(ctx, role, Email)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateActiveStatus(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserPassword(ctx context.Context, email string, newPasswordHash string) error {
	args := m.Called(ctx, email, newPasswordHash)
	return args.Error(0)
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context, page int, limit int) ([]domain.User, int64, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.User), args.Get(1).(int64), args.Error(2)
}

// MockUnactiveUserRepo mocks domain.UnactiveUserRepo
type MockUnactiveUserRepo struct {
	mock.Mock
}

func (m *MockUnactiveUserRepo) CreateUnactiveUser(ctx context.Context, user *domain.UnactivatedUser) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUnactiveUserRepo) FindByEmailUnactive(ctx context.Context, email string) (*domain.UnactivatedUser, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UnactivatedUser), args.Error(1)
}

func (m *MockUnactiveUserRepo) DeleteUnactiveUser(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockUnactiveUserRepo) UpdateActiveToken(ctx context.Context, email string, token string, expiry time.Time) error {
	args := m.Called(ctx, email, token, expiry)
	return args.Error(0)
}

// MockTokenRepository mocks domain.TokenRepository
type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) StoreRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepository) FindRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *MockTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepository) DeleteAllForUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// MockPasswordResetRepository mocks domain.PasswordResetRepository
type MockPasswordResetRepository struct {
	mock.Mock
}

func (m *MockPasswordResetRepository) Create(ctx context.Context, pr *domain.PasswordResetToken) error {
	args := m.Called(ctx, pr)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) GetByToken(ctx context.Context, token string) (*domain.PasswordResetToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PasswordResetToken), args.Error(1)
}

func (m *MockPasswordResetRepository) Delete(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func TestAuthUsecase_Register(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockUnactiveRepo := new(MockUnactiveUserRepo)

	// Set a dummy SMTP_PORT to prevent NewEmailService from returning nil
	t.Setenv("SMTP_PORT", "587")
	emailService := email.NewEmailService()

	uc := NewAuthUsecase(mockUserRepo, nil, nil, mockUnactiveRepo, emailService, nil)

	ctx := context.Background()
	emailAddr := "test@example.com"
	username := "testuser"
	password := "password"

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("FindByEmail", ctx, emailAddr).Return(nil, errors.New("not found")).Once()
		mockUnactiveRepo.On("FindByEmailUnactive", ctx, emailAddr).Return(nil, errors.New("not found")).Once()
		mockUnactiveRepo.On("CreateUnactiveUser", ctx, mock.AnythingOfType("*domain.UnactivatedUser")).Return(nil).Once()

		err := uc.Register(ctx, emailAddr, username, password)
		assert.NoError(t, err)

		mockUserRepo.AssertExpectations(t)
		mockUnactiveRepo.AssertExpectations(t)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockUserRepo.On("FindByEmail", ctx, emailAddr).Return(&domain.User{}, nil).Once()
		err := uc.Register(ctx, emailAddr, username, password)
		assert.Error(t, err)
		assert.Equal(t, "user already exists", err.Error())
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("unactivated user already exists", func(t *testing.T) {
		mockUserRepo.On("FindByEmail", ctx, emailAddr).Return(nil, errors.New("not found")).Once()
		mockUnactiveRepo.On("FindByEmailUnactive", ctx, emailAddr).Return(&domain.UnactivatedUser{}, nil).Once()
		err := uc.Register(ctx, emailAddr, username, password)
		assert.Error(t, err)
		assert.Equal(t, "user already exists please activate your account", err.Error())
		mockUserRepo.AssertExpectations(t)
		mockUnactiveRepo.AssertExpectations(t)
	})
}

func TestAuthUsecase_Login(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockUnactiveRepo := new(MockUnactiveUserRepo)
	mockTokenRepo := new(MockTokenRepository)
	jwt := auth.NewJWT("secret", "refresh_secret", time.Hour, time.Hour)
	uc := NewAuthUsecase(mockUserRepo, mockTokenRepo, jwt, mockUnactiveRepo, nil, nil)

	ctx := context.Background()
	emailAddr := "test@example.com"
	password := "password"
	hashedPassword, _ := uc.hasher.HashPassword(password)
	user := &domain.User{ID: primitive.NewObjectID().Hex(), Email: emailAddr, Password: hashedPassword, Activated: true, Role: "user"}

	t.Run("success", func(t *testing.T) {
		mockUnactiveRepo.On("FindByEmailUnactive", ctx, emailAddr).Return(nil, errors.New("not found")).Once()
		mockUserRepo.On("FindByEmail", ctx, emailAddr).Return(user, nil).Once()
		mockTokenRepo.On("StoreRefreshToken", ctx, mock.AnythingOfType("*domain.RefreshToken")).Return(nil).Once()

		accessToken, refreshToken, status, loggedInUser, err := uc.Login(ctx, emailAddr, password)

		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		assert.Equal(t, int(time.Hour.Seconds()), status)
		assert.Equal(t, user, loggedInUser)

		mockUnactiveRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockTokenRepo.AssertExpectations(t)
	})

	t.Run("user not activated", func(t *testing.T) {
		mockUnactiveRepo.On("FindByEmailUnactive", ctx, emailAddr).Return(&domain.UnactivatedUser{}, nil).Once()
		_, _, _, _, err := uc.Login(ctx, emailAddr, password)
		assert.Error(t, err)
		assert.Equal(t, "user not activated, please check your email for activation link", err.Error())
		mockUnactiveRepo.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockUnactiveRepo.On("FindByEmailUnactive", ctx, emailAddr).Return(nil, errors.New("not found")).Once()
		mockUserRepo.On("FindByEmail", ctx, emailAddr).Return(nil, errors.New("not found")).Once()
		_, _, _, _, err := uc.Login(ctx, emailAddr, password)
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})
}
