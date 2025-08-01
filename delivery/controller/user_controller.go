package controller

import (
	"g3-g65-bsp/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type emailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type profileUpdateRequest struct {
	Email   string             `json:"email" binding:"required,email"`
	Profile domain.UserProfile `json:"profile" binding:"required"`
}

type UserController struct {
	userOperations domain.UserOperations
}

func NewUserController(uuc domain.UserOperations) *UserController {
	return &UserController{
		userOperations: uuc,
	}
}

func (uc *UserController) HandlePromote(c *gin.Context) {
	var req emailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}

	ctx := c.Request.Context()
	err := uc.userOperations.Promote(ctx, req.Email)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "user promoted successfully"})
}

func (uc *UserController) HandleDemote(c *gin.Context) {
	var req emailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}

	ctx := c.Request.Context()
	err := uc.userOperations.Demote(ctx, req.Email)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, err)
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "user demoted successfully"})
}

func (uc *UserController) HandleUpdateUser(c *gin.Context) {
	var req profileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	ctx := c.Request.Context()
	if err := uc.userOperations.ProfileUpdate(ctx, &req.Profile, req.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}
