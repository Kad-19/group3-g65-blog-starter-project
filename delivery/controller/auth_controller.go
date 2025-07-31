package controller

import (
	"g3-g65-bsp/domain"
	"g3-g65-bsp/infrastructure/auth"
	"g3-g65-bsp/usecase"
	"net/http"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authUsecase usecase.AuthUsecase
	jwt         *auth.JWT
}

func NewAuthController(uc usecase.AuthUsecase, jwt *auth.JWT) *AuthController {
	return &AuthController{authUsecase: uc, jwt: jwt}
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