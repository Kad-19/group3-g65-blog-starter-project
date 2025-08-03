package usecase

import (
	"context"
	"errors"
	"g3-g65-bsp/domain"
	"g3-g65-bsp/infrastructure/auth"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthUsecase struct {
	userRepo  domain.UserRepository
	tokenRepo domain.TokenRepository
	hasher    *auth.PasswordHasher
	jwt       *auth.JWT
}

func NewAuthUsecase(
	ur domain.UserRepository,
	tr domain.TokenRepository,
	jwt *auth.JWT,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:  ur,
		tokenRepo: tr,
		hasher:    &auth.PasswordHasher{},
		jwt:       jwt,
	}
}

func (uc *AuthUsecase) Register(ctx context.Context, email, username, password string) (*domain.User, error) {
	if _, err := uc.userRepo.FindByEmail(ctx, email); err == nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := uc.hasher.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		Email:        email,
		Password:     hashedPassword,
		Role:         "User",
		Activated:    false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, email, password string) (string, string, int, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", 0, errors.New("invalid credentials")
	}

	if !uc.hasher.CompareHashAndPassword(user.Password, password) {
		return "", "", 0, errors.New("invalid credentials")
	}

	accessToken, err := uc.jwt.GenerateAccessToken(user.ID.Hex(), user.Role)
	if err != nil {
		return "", "", 0, errors.New("failed to generate token")
	}

	refreshToken, err := uc.jwt.GenerateRefreshToken()
	if err != nil {
		return "", "", 0, errors.New("failed to generate refresh token")
	}

	if err := uc.tokenRepo.StoreRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(uc.jwt.RefreshExpiry)); err != nil {
		return "", "", 0, errors.New("failed to store refresh token")
	}

	return accessToken, refreshToken, int(uc.jwt.AccessExpiry.Seconds()), nil
}

func (uc *AuthUsecase) RefreshTokens(ctx context.Context, refreshToken string) (string, string, int, error) {
	// 1. Validate and delete the old refresh token
	userID, err := uc.tokenRepo.FindAndDeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", 0, errors.New("invalid refresh token")
	}

	// 2. Get user details
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", "", 0, errors.New("user not found")
	}

	// 3. Generate new tokens
	newAccessToken, err := uc.jwt.GenerateAccessToken(user.ID.Hex(), user.Role)
	if err != nil {
		return "", "", 0, errors.New("failed to generate access token")
	}

	newRefreshToken, err := uc.jwt.GenerateRefreshToken()
	if err != nil {
		return "", "", 0, errors.New("failed to generate refresh token")
	}

	// 4. Store new refresh token
	expiry := time.Now().Add(uc.jwt.RefreshExpiry)
	if err := uc.tokenRepo.StoreRefreshToken(ctx, user.ID, newRefreshToken, expiry); err != nil {
		return "", "", 0, errors.New("failed to store refresh token")
	}

	return newAccessToken, newRefreshToken, int(uc.jwt.AccessExpiry.Seconds()), nil
}

// Logout (single device)
func (uc *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	_, err := uc.tokenRepo.FindAndDeleteRefreshToken(ctx, refreshToken)
	return err
}

// LogoutAll (all devices)
func (uc *AuthUsecase) LogoutAll(ctx context.Context, userID primitive.ObjectID) error {
	return uc.tokenRepo.DeleteAllForUser(ctx, userID)
}