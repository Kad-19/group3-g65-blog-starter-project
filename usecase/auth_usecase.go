package usecase

import (
	"context"
	"errors"
	"fmt"
	"g3-g65-bsp/domain"
	"g3-g65-bsp/infrastructure/auth"
	"g3-g65-bsp/infrastructure/email"
	"g3-g65-bsp/utils"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthUsecase struct {
	userRepo       domain.UserRepository
	tokenRepo      domain.TokenRepository
	hasher         *auth.PasswordHasher
	jwt            *auth.JWT
	activationRepo domain.ActivationTokenRepository
	emailService   *email.EmailService
	passRepo       domain.PasswordResetRepository
}

func NewAuthUsecase(
	ur domain.UserRepository,
	tr domain.TokenRepository,
	jwt *auth.JWT,
	ar domain.ActivationTokenRepository,
	es *email.EmailService,
	passRepo domain.PasswordResetRepository,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:       ur,
		tokenRepo:      tr,
		hasher:         &auth.PasswordHasher{},
		jwt:            jwt,
		activationRepo: ar,
		emailService:   es,
		passRepo:       passRepo,
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

	activation_token, err := utils.CreateActivationToken(user.Email, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	if err := uc.activationRepo.Create(ctx, activation_token); err != nil {
		return nil, err
	}

	activationLink := ""
	// activationLink := "" + activation_token
	go func() { // Send email asynchronously
		err := uc.emailService.SendActivationEmail(user.Email, activationLink)
		if err != nil {
			// Log the error, but don't fail the registration process
			fmt.Printf("Failed to send activation email: %v\n", err)
		}
	}()
	
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

func (uc *AuthUsecase) ActivateUser(ctx context.Context, token string) error {
	activate_token, err := uc.activationRepo.GetByToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if time.Now().After(activate_token.ExpiresAt) {
		uc.activationRepo.Delete(ctx, token) // Clean up expired token
		return errors.New("token has expired")
	}

	user, err := uc.userRepo.FindByEmail(ctx, activate_token.Email)
	if err != nil {
		return errors.New("invalid token")
	}

	if user.Activated {
		return errors.New("account is already active")
	}

	err = uc.userRepo.UpdateActiveStatus(ctx, user.Email)
	if err != nil {
		return err
	}

	uc.activationRepo.Delete(ctx, token)
	return nil
}

func (uc *AuthUsecase) InitiateResetPassword(ctx context.Context, email string) error {
	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("if an account exists for this email, a password reset link will be sent")
	}

	password_token, _ := utils.CreateResetToken(email, 1*time.Hour)
	if err := uc.passRepo.Create(ctx, password_token); err != nil {
		return err
	}

	// resetLink := "http://your-frontend-domain.com/reset-password?token=" + tokenValue
	resetLink := ""
	go func() {
		err := uc.emailService.SendPasswordResetEmail(user.Email, resetLink)
		if err != nil {
			// Log the error.
			// fmt.Printf("Failed to send password reset email: %v\n", err)
		}
	}()

	return errors.New("if an account exists for this email, a password reset link will be sent")
}

func (uc *AuthUsecase) ResetPassword(c context.Context, token, newPassword string) error {
	resetToken, err := uc.passRepo.GetByToken(c, token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	if time.Now().After(resetToken.ExpiresAt) {
		uc.passRepo.Delete(c, token) // Clean up expired token.
		return errors.New("token has expired")
	}

	user, err := uc.userRepo.FindByEmail(c, resetToken.Email)
	if err != nil {
		return errors.New("user not found")
	}

	hashedPassword, err := uc.hasher.HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := uc.userRepo.UpdateUserPassword(c, user.Email, hashedPassword); err != nil {
		return err
	}

	uc.passRepo.Delete(c, token)
	return nil
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
