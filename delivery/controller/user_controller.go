package controller

import (
	"context"
	"g3-g65-bsp/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userOperations domain.UserUseCase
}

type emailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func NewUserController(uuc domain.UserUseCase) *UserController {
	return &UserController{
		userOperations: uuc,
	}
}

func (uc *UserController) ChangeUserRole(c *gin.Context, roleChange func(context.Context, string) error, successMessage string) {
	var req emailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	if err := roleChange(ctx, req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": successMessage})
}
func (uc *UserController) HandlePromote(c *gin.Context) {
	uc.ChangeUserRole(c, uc.userOperations.Demote, "user promoted successfully")
}

func (uc *UserController) HandleDemote(c *gin.Context) {
	uc.ChangeUserRole(c, uc.userOperations.Demote, "user demoted successfully")
}

func (uc *UserController) HandleUpdateUser(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ObjectID"})
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
	if err := uc.userOperations.ProfileUpdate(ctx, userID, bio, contactinfo, fileReader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
