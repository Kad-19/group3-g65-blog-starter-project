package usecase

import (
	"context"
	"errors"
	"time"
	"g3-g65-bsp/domain"
	"g3-g65-bsp/infrastructure/auth"
	"g3-g65-bsp/repository"
)

type AuthUsecase struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
	hasher    *auth.PasswordHasher
	jwt       *auth.JWT
}

func NewAuthUsecase(
	ur repository.UserRepository,
	tr repository.TokenRepository,
	jwt *auth.JWT,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:  ur,
		tokenRepo: tr,
		hasher:    &auth.PasswordHasher{},
		jwt:       jwt,
	}
}

func (uc *AuthUsecase) Register(ctx context.Context, req domain.UserCreateRequest) (*domain.UserResponse, error) {
	if _, err := uc.userRepo.FindByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := uc.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		Password:     hashedPassword,
		Role:         "User",
		Activated:    false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (uc *AuthUsecase) Login(ctx context.Context, email, password string) (string, string, int, error) {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", 0, errors.New("invalid credentials")
	}

	if !uc.hasher.CompareHashAndPassword(password, user.Password) {
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