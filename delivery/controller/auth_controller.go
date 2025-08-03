package controller

import (
	"fmt"
	"g3-g65-bsp/domain"
	"g3-g65-bsp/infrastructure/auth"
	"net/http"
	"time"

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

// UserResponse represents safe user data for API responses (DTO)
type UserResponse struct {
	ID        string      `json:"id"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	Role      string      `json:"role"`
	Profile   domain.UserProfile `json:"profile"`
	CreatedAt time.Time   `json:"created_at"`
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

	user, err := c.authUsecase.Register(ctx, req.Email, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, UserResponse{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Profile:   user.Profile,
		CreatedAt: user.CreatedAt,
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req UserLoginRequest

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
	fmt.Println("User ID from context:", userID)
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