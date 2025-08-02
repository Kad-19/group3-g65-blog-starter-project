package controller

import (
	"g3-g65-bsp/domain"
	"g3-g65-bsp/infrastructure/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUsecase domain.AuthUsecase
	jwt         *auth.JWT
}

func NewAuthController(uc domain.AuthUsecase, jwt *auth.JWT) *AuthController {
	return &AuthController{authUsecase: uc, jwt: jwt}
}

type ForgotEmail struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req domain.UserCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.authUsecase.Register(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req domain.UserLoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, expiresIn, err := c.authUsecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    expiresIn,
	})
}

func (c *AuthController) ActivateUser(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	err := c.authUsecase.ActivateUser(ctx.Request.Context(), token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "account activated successfully"})
}

func (ac *AuthController) ForgotPassWord(ctx *gin.Context) {
	var req ForgotEmail
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ac.authUsecase.InitiateResetPassword(ctx.Request.Context(), req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "if an account exists, a password reset link has been sent"})
}

func (ac *AuthController) Reset(c *gin.Context) {
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
