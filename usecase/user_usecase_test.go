package usecase

import (
	"context"
	"g3-g65-bsp/domain"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockImageUploader is a mock implementation of the ImageUploader interface.
type MockImageUploader struct {
	mock.Mock
}

func (m *MockImageUploader) UploadImage(ctx context.Context, file io.Reader, folder string) (string, error) {
	args := m.Called(ctx, file, folder)
	return args.String(0), args.Error(1)
}

func TestUserUsecase_Promote(t *testing.T) {
	ctx := context.Background()
	adminID := "admin123"
	userID := "user123"
	userEmail := "user@example.com"

	t.Run("successful promotion", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		uc := NewUserUsecase(mockUserRepo, nil)
		user := &domain.User{ID: userID, Email: userEmail, Role: string(domain.RoleUser)}
		mockUserRepo.On("FindByEmail", ctx, userEmail).Return(user, nil).Once()
		mockUserRepo.On("UpdateUserRole", ctx, string(domain.RoleAdmin), userEmail).Return(nil).Once()

		err := uc.Promote(ctx, adminID, userEmail)
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("cannot promote self", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		uc := NewUserUsecase(mockUserRepo, nil)
		user := &domain.User{ID: userID, Email: userEmail, Role: string(domain.RoleUser)}
		mockUserRepo.On("FindByEmail", ctx, userEmail).Return(user, nil).Once()

		err := uc.Promote(ctx, userID, userEmail)
		assert.Error(t, err)
		assert.Equal(t, ErrSelfRoleChange, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("user already has the target role", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		uc := NewUserUsecase(mockUserRepo, nil)
		adminUser := &domain.User{ID: "admin456", Email: "admin@example.com", Role: string(domain.RoleAdmin)}
		mockUserRepo.On("FindByEmail", ctx, "admin@example.com").Return(adminUser, nil).Once()

		err := uc.Promote(ctx, adminID, "admin@example.com")
		assert.Error(t, err)
		assert.Equal(t, ErrAlreadyHasRole, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_Demote(t *testing.T) {
	ctx := context.Background()
	adminID := "admin123"
	userID := "user123"
	userEmail := "user@example.com"

	t.Run("successful demotion", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		uc := NewUserUsecase(mockUserRepo, nil)
		user := &domain.User{ID: userID, Email: userEmail, Role: string(domain.RoleAdmin)}
		mockUserRepo.On("FindByEmail", ctx, userEmail).Return(user, nil).Once()
		mockUserRepo.On("UpdateUserRole", ctx, string(domain.RoleUser), userEmail).Return(nil).Once()

		err := uc.Demote(ctx, adminID, userEmail)
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("cannot demote self", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		uc := NewUserUsecase(mockUserRepo, nil)
		user := &domain.User{ID: userID, Email: userEmail, Role: string(domain.RoleAdmin)}
		mockUserRepo.On("FindByEmail", ctx, userEmail).Return(user, nil).Once()

		err := uc.Demote(ctx, userID, userEmail)
		assert.Error(t, err)
		assert.Equal(t, ErrSelfRoleChange, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("user already has the target role", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		uc := NewUserUsecase(mockUserRepo, nil)
		normalUser := &domain.User{ID: "user456", Email: "normal@example.com", Role: string(domain.RoleUser)}
		mockUserRepo.On("FindByEmail", ctx, "normal@example.com").Return(normalUser, nil).Once()

		err := uc.Demote(ctx, adminID, "normal@example.com")
		assert.Error(t, err)
		assert.Equal(t, ErrAlreadyHasRole, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_ProfileUpdate(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockImageUploader := new(MockImageUploader)
	uc := NewUserUsecase(mockUserRepo, mockImageUploader)

	ctx := context.Background()
	userID := "user123"
	userEmail := "user@example.com"
	user := &domain.User{ID: userID, Email: userEmail}
	bio := "New bio"
	contactInfo := "new-contact"
	imageURL := "http://example.com/new-image.jpg"
	file := strings.NewReader("fake image data")

	mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
	mockImageUploader.On("UploadImage", ctx, file, "profile").Return(imageURL, nil).Once()
	mockUserRepo.On("UpdateUserProfile", ctx, bio, contactInfo, imageURL, userEmail).Return(nil).Once()

	err := uc.ProfileUpdate(ctx, userID, bio, contactInfo, file)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockImageUploader.AssertExpectations(t)
}

func TestUserUsecase_GetAllUsers(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	uc := NewUserUsecase(mockUserRepo, nil)

	ctx := context.Background()
	page := 1
	limit := 10
	users := []domain.User{{ID: "user1"}, {ID: "user2"}}
	total := int64(2)

	mockUserRepo.On("GetAllUsers", ctx, page, limit).Return(users, total, nil).Once()

	resultUsers, resultTotal, err := uc.GetAllUsers(ctx, page, limit)

	assert.NoError(t, err)
	assert.Equal(t, users, resultUsers)
	assert.Equal(t, total, resultTotal)
	mockUserRepo.AssertExpectations(t)
}
