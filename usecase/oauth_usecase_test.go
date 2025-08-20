package usecase

import (
	"context"
	"g3-g65-bsp/domain"
	"g3-g65-bsp/infrastructure/auth"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockOAuthUserRepository is a mock implementation of the UserRepository for OAuth tests.
type MockOAuthUserRepository struct {
	mock.Mock
}

func (m *MockOAuthUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockOAuthUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockOAuthUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockOAuthUserRepository) UpdateUserProfile(ctx context.Context, bio string, contactInfo string, imagePath string, Email string) error {
	args := m.Called(ctx, bio, contactInfo, imagePath, Email)
	return args.Error(0)
}

func (m *MockOAuthUserRepository) UpdateUserRole(ctx context.Context, role string, Email string) error {
	args := m.Called(ctx, role, Email)
	return args.Error(0)
}

func (m *MockOAuthUserRepository) UpdateActiveStatus(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockOAuthUserRepository) UpdateUserPassword(ctx context.Context, email string, newPasswordHash string) error {
	args := m.Called(ctx, email, newPasswordHash)
	return args.Error(0)
}

func (m *MockOAuthUserRepository) GetAllUsers(ctx context.Context, page int, limit int) ([]domain.User, int64, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]domain.User), args.Get(1).(int64), args.Error(2)
}

// MockOAuthTokenRepository is a mock implementation of the TokenRepository for OAuth tests.
type MockOAuthTokenRepository struct {
	mock.Mock
}

func (m *MockOAuthTokenRepository) StoreRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockOAuthTokenRepository) FindRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RefreshToken), args.Error(1)
}

func (m *MockOAuthTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockOAuthTokenRepository) DeleteAllForUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestOAuthUsecase_Placeholder(t *testing.T) {
	// This is a placeholder test to make the file compile.
	// Testing the full OAuthLogin method is complex as it involves external calls.
	mockUserRepo := new(MockOAuthUserRepository)
	mockTokenRepo := new(MockOAuthTokenRepository)
	jwtService := auth.NewJWT("test-secret", "test-refresh-secret", 1*time.Hour, 24*time.Hour)

	_ = NewOAuthUsecase(mockUserRepo, mockTokenRepo, jwtService)

}
