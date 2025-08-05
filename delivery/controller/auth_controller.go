package controller

import (
	"g3-g65-bsp/domain"
	"g3-g65-bsp/infrastructure/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type AuthController struct {
	authUsecase domain.AuthUsecase
	jwt         *auth.JWT
}

func NewAuthController(uc domain.AuthUsecase, jwt *auth.JWT) *AuthController {
	return &AuthController{authUsecase: uc, jwt: jwt}
}

type ActivatEmail struct {
	Email string `json:"email" binding:"required,email"`
}

type ForgotEmail struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
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
		"user":          user,
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

	ctx.JSON(http.StatusOK, gin.H{"message": "account activated successfully"})
}

func (ac *AuthController) ResendActivationEmail(ctx *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

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
	var req ForgotEmail
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.authUsecase.ForgotPassword(ctx.Request.Context(), req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "if an account exists, a password reset link has been sent"})
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
func (c *AuthController) Refresh(ctx *gin.Context) {
	if ctx.GetHeader("Authorization") == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, expiresIn, err := c.authUsecase.RefreshTokens(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
	})
}

// Logout (single device)
func (c *AuthController) Logout(ctx *gin.Context) {
	if ctx.GetHeader("Authorization") == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

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
	if ctx.GetHeader("Authorization") == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := c.authUsecase.LogoutAll(ctx.Request.Context(), objID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "logout failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logged out from all devices"})
}

func (ac *AuthController) ResendActivationToken(c *gin.Context) {
	var req ActivatEmail
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := ac.authUsecase.Reactivate(c.Request.Context(), req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "check your email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successfully sent reactivation"})
}
