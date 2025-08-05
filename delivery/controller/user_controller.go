package controller

import (
	"context"
	"g3-g65-bsp/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userUsecase domain.UserUsecase
}

type emailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func NewUserController(uuc domain.UserUsecase) *UserController {
	return &UserController{
		userUsecase: uuc,
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
	err := uc.userUsecase.Promote(ctx, req.Email)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": successMessage})
}
func (uc *UserController) HandlePromote(c *gin.Context) {
	uc.ChangeUserRole(c, uc.userOperations.Demote, "user promoted successfully")
}

func (uc *UserController) HandleDemote(c *gin.Context) {
	uc.ChangeUserRole(c, uc.userOperations.Demote, "user demoted successfully")
	var req emailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}

	ctx := c.Request.Context()
	err := uc.userUsecase.Demote(ctx, req.Email)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "user demoted successfully"})
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
	if err := uc.userUsecase.ProfileUpdate(ctx, userID, bio, contactinfo, fileReader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
