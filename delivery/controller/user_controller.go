package controller

import (
	"context"
	"g3-g65-bsp/domain"
	"net/http"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUsecase domain.UserUsecase
}

func NewUserController(uuc domain.UserUsecase) *UserController {
	return &UserController{
		userUsecase: uuc,
	}
}

func (uc *UserController) ChangeUserRole(c *gin.Context, roleChange func(context.Context, string) error, successMessage string) {
	var req EmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	if err := roleChange(ctx, req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		if err := roleChange(ctx, req.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}
func (uc *UserController) HandlePromote(c *gin.Context) {
	uc.ChangeUserRole(c, uc.userUsecase.Promote, "user promoted successfully")
}

func (uc *UserController) HandleDemote(c *gin.Context) {
	uc.ChangeUserRole(c, uc.userUsecase.Demote, "user demoted successfully")
}

func (uc *UserController) HandleUpdateUser(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	bio := c.PostForm("bio")
	contactinfo := c.PostForm("contact_info")

	file, err := c.FormFile("profile_picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileReader.Close()

	ctx := c.Request.Context()
	if err := uc.userUsecase.ProfileUpdate(ctx, userID, bio, contactinfo, fileReader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
