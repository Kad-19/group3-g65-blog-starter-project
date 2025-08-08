package controller

import (
	"g3-g65-bsp/domain"
	"g3-g65-bsp/infrastructure/auth"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// UserDTO represents the user data transfer object
type UserDTO struct {
	ID        string         `json:"id"`
	Username  string         `json:"username" validate:"required,min=3,max=50"`
	Email     string         `json:"email" validate:"required,email"`
	Password  string         `json:"-"`
	Role      string         `json:"role"`
	Activated bool           `json:"activated"`
	Profile   UserProfileDTO `json:"profile"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// UserProfileDTO represents the user profile data transfer object
type UserProfileDTO struct {
	Bio               string `json:"bio" validate:"min=3,max=500"`
	ProfilePictureURL string `json:"profile_picture_url"`
	ContactInfo       string `json:"contact_information" validate:"max=100"`
}

// UnactivatedUserDTO represents a user who has not yet activated their account
type UnactivatedUserDTO struct {
	ID                    string     `bson:"_id,omitempty" json:"id"`
	Username              string     `bson:"username" json:"username" validate:"required,min=3,max=50"`
	Email                 string     `bson:"email" json:"email" validate:"required,email"`
	Password              string     `bson:"password" json:"-"`
	Activated             bool       `bson:"activated" json:"activated"`
	ActivationToken       string     `bson:"activation_token,omitempty" json:"activation_token,omitempty"`
	ActivationTokenExpiry *time.Time `bson:"activation_token_expiry,omitempty" json:"activation_token_expiry,omitempty"`
	CreatedAt             time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt             time.Time  `bson:"updated_at" json:"updated_at"`
}

// ConvertToDomain converts UserDTO to domain.User
func (dto *UserDTO) ConvertToUserDomain() *domain.User {
	return &domain.User{
		ID:        dto.ID,
		Username:  dto.Username,
		Email:     dto.Email,
		Password:  dto.Password,
		Role:      dto.Role,
		Activated: dto.Activated,
		Profile: domain.UserProfile{
			Bio:               dto.Profile.Bio,
			ProfilePictureURL: dto.Profile.ProfilePictureURL,
			ContactInfo:       dto.Profile.ContactInfo,
		},
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

// ConvertToDomain converts UnactivatedUserDTO to domain.UnactivatedUser
func (dto *UnactivatedUserDTO) ConvertToUnactivatedUserDomain() *domain.UnactivatedUser {
	return &domain.UnactivatedUser{
		ID:                    dto.ID,
		Username:              dto.Username,
		Email:                 dto.Email,
		Password:              dto.Password,
		Activated:             dto.Activated,
		ActivationToken:       dto.ActivationToken,
		ActivationTokenExpiry: dto.ActivationTokenExpiry,
		CreatedAt:             dto.CreatedAt,
		UpdatedAt:             dto.UpdatedAt,
	}
}

// ConvertToDTO converts domain.User to UserDTO
func ConvertToUserDTO(u *domain.User) *UserDTO {
	return &UserDTO{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		Activated: u.Activated,
		Profile: UserProfileDTO{
			Bio:               u.Profile.Bio,
			ProfilePictureURL: u.Profile.ProfilePictureURL,
			ContactInfo:       u.Profile.ContactInfo,
		},
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ConvertToDTO converts domain.UnactivatedUser to UnactivatedUserDTO
func ConvertToUnactivatedUserDTO(u *domain.UnactivatedUser) *UnactivatedUserDTO {
	return &UnactivatedUserDTO{
		ID:                    u.ID,
		Username:              u.Username,
		Email:                 u.Email,
		Password:              u.Password,
		Activated:             u.Activated,
		ActivationToken:       u.ActivationToken,
		ActivationTokenExpiry: u.ActivationTokenExpiry,
		CreatedAt:             u.CreatedAt,
		UpdatedAt:             u.UpdatedAt,
	}
}

// UserCreateRequest represents registration payload (DTO)
type UserCreateRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// UserLoginRequest represents login payload (DTO)
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type EmailReq struct {
	Email string `json:"email" binding:"required,email"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type PasswordResetRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type AuthController struct {
	authUsecase domain.AuthUsecase
	jwt         *auth.JWT
}

func NewAuthController(uc domain.AuthUsecase, jwt *auth.JWT) *AuthController {
	return &AuthController{authUsecase: uc, jwt: jwt}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.authUsecase.Register(ctx, req.Email, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "user registered successfully please check your email to activate your account"})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req UserLoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, expiresIn, user, err := c.authUsecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
		"user":          ConvertToUserDTO(user),
	})
}

func (c *AuthController) ActivateUser(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "activation token is required"})
		return
	}

	email := ctx.Query("email")
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	err := c.authUsecase.ActivateUser(ctx.Request.Context(), token, email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.HTML(http.StatusOK, "activation_success.html", gin.H{
		"email": email,
	})

}

func (ac *AuthController) ResendActivationEmail(ctx *gin.Context) {
	var req EmailReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.authUsecase.ResendActivationEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "activation email resent successfully"})
}

func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var req EmailReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.authUsecase.ForgotPassword(ctx.Request.Context(), req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "reset is sent to email successfully"})
}

func (ac *AuthController) ResetPassword(c *gin.Context) {
	var request PasswordResetRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.authUsecase.ResetPassword(c.Request.Context(), request.Token, request.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password has been reset successfully"})
}

// Refresh Tokens
func (c *AuthController) RefreshAccessToken(ctx *gin.Context) {
	// Accept refresh token from either header or JSON body
	var req RefreshTokenRequest

	refreshToken := ctx.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "refresh token is required"})
			return
		}
		refreshToken = req.RefreshToken
	}

	accessToken, refreshTokenNew, expiresIn, err := c.authUsecase.RefreshTokens(ctx.Request.Context(), refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshTokenNew,
		"expires_in":    expiresIn,
	})
}

// Logout (single device)
func (c *AuthController) Logout(ctx *gin.Context) {
	var req RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.authUsecase.Logout(ctx.Request.Context(), req.RefreshToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}

// LogoutAll (all devices)
func (c *AuthController) LogoutAll(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	objID := userID.(string)
	if objID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	if err := c.authUsecase.LogoutAll(ctx.Request.Context(), objID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logged out from all devices"})
}
